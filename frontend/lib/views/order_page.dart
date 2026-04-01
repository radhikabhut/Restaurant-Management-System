import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/order_provider.dart';
import '../providers/auth_provider.dart';
import '../models/order.dart';
import 'create_order_page.dart';

class OrderPage extends ConsumerWidget {
  const OrderPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authProvider);
    final role = authState.role ?? 'waiter';
    final orderList = ref.watch(orderListProvider((userId: null, status: null)));

    return Scaffold(
      body: orderList.when(
        data: (orders) => orders.isEmpty
            ? const Center(child: Text('No orders yet'))
            : ListView.builder(
                padding: const EdgeInsets.all(16),
                itemCount: orders.length,
                itemBuilder: (context, index) {
                  final order = orders[index];
                  return Card(
                    margin: const EdgeInsets.only(bottom: 12),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(16),
                    ),
                    child: ExpansionTile(
                      leading: order.isDeleted
                          ? const Icon(Icons.delete_forever, color: Colors.grey)
                          : null,
                      title: Text(
                        'Table #${order.tableNumber} - Order #${order.id.split('-').first}',
                        style: TextStyle(
                          fontWeight: FontWeight.bold,
                          color: order.isDeleted ? Colors.grey : null,
                          decoration: order.isDeleted ? TextDecoration.lineThrough : null,
                        ),
                      ),
                      subtitle: Text(
                        'Total: \$${order.totalAmount} • ${order.isDeleted ? "DELETED" : order.status}',
                        style: TextStyle(color: order.isDeleted ? Colors.grey : null),
                      ),
                      children: [
                        ...order.items.map((it) => ListTile(
                              title: Text(it.menuName,
                                  style: TextStyle(
                                    fontWeight: FontWeight.bold,
                                    color: order.isDeleted ? Colors.grey : null,
                                  )),
                              subtitle: Text('Qty: ${it.quantity} • \$${it.price}',
                                  style: TextStyle(color: order.isDeleted ? Colors.grey : null)),
                              trailing: !order.isDeleted && (role == 'admin' || role == 'waiter')
                                  ? IconButton(
                                      icon: const Icon(Icons.remove_circle_outline, color: Colors.red, size: 20),
                                      onPressed: () => _showDeleteItemConfirmation(context, ref, order.id, it),
                                    )
                                  : null,
                            )),
                        const Divider(),
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                children: [
                                  if (!order.isDeleted && role == 'kitchen' && order.status != 'completed' && order.status != 'cancelled')
                                    Wrap(
                                      spacing: 8,
                                      children: [
                                        if (order.status == 'pending')
                                          _StatusButton(
                                            label: 'Prepare',
                                            color: Colors.blue,
                                            onPressed: () => ref
                                                .read(orderActionProvider.notifier)
                                                .updateOrder(order.id, status: 'preparing'),
                                          ),
                                        if (order.status == 'preparing')
                                          _StatusButton(
                                            label: 'Complete',
                                            color: Colors.green,
                                            onPressed: () => ref
                                                .read(orderActionProvider.notifier)
                                                .updateOrder(order.id, status: 'completed'),
                                          ),
                                        _StatusButton(
                                          label: 'Cancel',
                                          color: Colors.red,
                                          onPressed: () => ref
                                              .read(orderActionProvider.notifier)
                                              .updateOrder(order.id, status: 'cancelled'),
                                        ),
                                      ],
                                    )
                                  else
                                    Text(
                                      order.isDeleted ? 'Status: DELETED' : 'Status: ${order.status.toUpperCase()}',
                                      style: TextStyle(
                                        fontWeight: FontWeight.bold,
                                        color: order.isDeleted ? Colors.grey : null,
                                      ),
                                    ),
                                  Row(
                                    children: [
                                      if (!order.isDeleted && (role == 'admin' || role == 'waiter'))
                                        Row(
                                          children: [
                                            _StatusButton(
                                              label: 'Add Items',
                                              color: Colors.orange,
                                              onPressed: () => Navigator.push(
                                                context,
                                                MaterialPageRoute(
                                                  builder: (context) => CreateOrderPage(
                                                    orderId: order.id,
                                                    existingTableNumber: order.tableNumber,
                                                  ),
                                                ),
                                              ),
                                            ),
                                            const SizedBox(width: 8),
                                            _StatusButton(
                                              label: 'Generate Bill',
                                              color: Colors.purple,
                                              onPressed: () => _showBillDialog(context, ref, order),
                                            ),
                                          ],
                                        ),
                                      if (!order.isDeleted && (role == 'admin' || role == 'waiter'))
                                        IconButton(
                                          icon: const Icon(Icons.delete_outline, color: Colors.red),
                                          onPressed: () => _showDeleteConfirmation(context, ref, order),
                                        ),
                                    ],
                                  ),
                                ],
                              ),
                            ],
                          ),
                        )
                      ],
                    ),
                  );
                },
              ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (err, stack) => Center(child: Text('Error: $err')),
      ),
      floatingActionButton: role != 'kitchen'
          ? FloatingActionButton(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => const CreateOrderPage()),
                );
              },
              child: const Icon(Icons.add),
            )
          : null,
    );
  }

  void _showBillDialog(BuildContext context, WidgetRef ref, Order order) async {
    showDialog(
      context: context,
      builder: (context) => const Center(child: CircularProgressIndicator()),
    );

    try {
      final bill = await ref.read(orderActionProvider.notifier).generateBill(order.id);
      if (context.mounted) {
        Navigator.pop(context); // Close loading
        showDialog(
          context: context,
          builder: (context) => AlertDialog(
            title: const Text('Bill Details', style: TextStyle(fontWeight: FontWeight.bold)),
            content: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('Table: ${bill['tableNumber']}', style: const TextStyle(fontWeight: FontWeight.bold)),
                  const SizedBox(height: 8),
                  ... (bill['items'] as List).map((it) => Padding(
                    padding: const EdgeInsets.symmetric(vertical: 2.0),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Expanded(child: Text('${it['menuName']} x${it['quantity']}')),
                        Text('\$${(it['subtotal'] as num).toStringAsFixed(2)}'),
                      ],
                    ),
                  )),
                  const Divider(),
                  _BillRow(label: 'Subtotal', value: bill['totalAmount']),
                  _BillRow(label: 'Tax (5%)', value: bill['tax']),
                  _BillRow(label: 'Service Charge (10%)', value: bill['serviceCharge']),
                  const Divider(),
                  _BillRow(label: 'Grand Total', value: bill['grandTotal'], isBold: true),
                ],
              ),
            ),
            actions: [
              TextButton(onPressed: () => Navigator.pop(context), child: const Text('Close')),
              ElevatedButton(
                onPressed: () {
                  // Print logic could go here
                  Navigator.pop(context);
                },
                child: const Text('Print'),
              ),
            ],
          ),
        );
      }
    } catch (e) {
      if (context.mounted) {
        Navigator.pop(context);
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    }
  }

  void _showDeleteItemConfirmation(BuildContext context, WidgetRef ref, String orderId, OrderItem item) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Remove Item'),
        content: Text('Are you sure you want to remove ${item.menuName} from this order?'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context), child: const Text('Cancel')),
          TextButton(
            onPressed: () async {
              await ref.read(orderActionProvider.notifier).updateOrder(orderId, deletedItemIds: [item.id]);
              if (context.mounted) Navigator.pop(context);
            },
            child: const Text('Remove', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  }

  void _showDeleteConfirmation(BuildContext context, WidgetRef ref, Order order) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Order'),
        content: Text('Are you sure you want to delete order #${order.id.split('-').first}?'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context), child: const Text('Cancel')),
          TextButton(
            onPressed: () async {
              final userId = ref.read(authProvider).userId ?? "00000000-0000-0000-0000-000000000000";
              await ref.read(orderActionProvider.notifier).deleteOrder(order.id, userId);
              if (context.mounted) Navigator.pop(context);
            },
            child: const Text('Delete', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  }
}

class _StatusButton extends StatelessWidget {
  final String label;
  final Color color;
  final VoidCallback onPressed;

  const _StatusButton({
    required this.label,
    required this.color,
    required this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      onPressed: onPressed,
      style: ElevatedButton.styleFrom(
        backgroundColor: color.withAlpha(50),
        foregroundColor: color,
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        minimumSize: Size.zero,
        tapTargetSize: MaterialTapTargetSize.shrinkWrap,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      ),
      child: Text(label, style: const TextStyle(fontSize: 12, fontWeight: FontWeight.bold)),
    );
  }
}

class _BillRow extends StatelessWidget {
  final String label;
  final dynamic value;
  final bool isBold;

  const _BillRow({required this.label, required this.value, this.isBold = false});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: TextStyle(fontWeight: isBold ? FontWeight.bold : FontWeight.normal)),
          Text(
            '\$${(value as num).toStringAsFixed(2)}',
            style: TextStyle(fontWeight: isBold ? FontWeight.bold : FontWeight.normal, fontSize: isBold ? 18 : 14),
          ),
        ],
      ),
    );
  }
}
