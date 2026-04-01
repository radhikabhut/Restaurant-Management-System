import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/menu_provider.dart';
import '../providers/order_provider.dart';
import '../models/order.dart';
import '../models/menu_item.dart';
import '../providers/auth_provider.dart';

class CreateOrderPage extends ConsumerStatefulWidget {
  final String? orderId;
  final String? existingTableNumber;

  const CreateOrderPage({super.key, this.orderId, this.existingTableNumber});

  @override
  ConsumerState<CreateOrderPage> createState() => _CreateOrderPageState();
}

class _CreateOrderPageState extends ConsumerState<CreateOrderPage> {
  final Map<String, int> _selectedItems = {}; // menuId -> quantity
  final TextEditingController _searchController = TextEditingController();
  final TextEditingController _tableNumberController = TextEditingController(); // Added
  String _searchQuery = '';

  @override
  void initState() {
    super.initState();
    if (widget.existingTableNumber != null) {
      _tableNumberController.text = widget.existingTableNumber!;
    }
  }

  @override
  void dispose() {
    _searchController.dispose();
    _tableNumberController.dispose(); // Added
    super.dispose();
  }

  void _addItem(MenuItem item) {
    setState(() {
      _selectedItems[item.id] = (_selectedItems[item.id] ?? 0) + 1;
    });
  }

  void _removeItem(MenuItem item) {
    setState(() {
      if (_selectedItems.containsKey(item.id)) {
        if (_selectedItems[item.id]! > 1) {
          _selectedItems[item.id] = _selectedItems[item.id]! - 1;
        } else {
          _selectedItems.remove(item.id);
        }
      }
    });
  }

  double _calculateTotal(List<MenuItem> menus) {
    double total = 0;
    _selectedItems.forEach((menuId, qty) {
      final menu = menus.firstWhere((m) => m.id == menuId);
      total += menu.price * qty;
    });
    return total;
  }

  @override
  Widget build(BuildContext context) {
    final menuList = ref.watch(menuListProvider);
    final orderAction = ref.watch(orderActionProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('New Order'),
      ),
      body: menuList.when(
        data: (menus) {
          final filteredMenus = menus.where((m) {
            final matchesSearch = m.name.toLowerCase().contains(_searchQuery.toLowerCase()) ||
                m.category.toLowerCase().contains(_searchQuery.toLowerCase());
            return m.isAvailable && matchesSearch;
          }).toList();

          return Column(
            children: [
              if (widget.orderId == null)
                Padding(
                  padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
                  child: TextField(
                    controller: _tableNumberController,
                    decoration: InputDecoration(
                      hintText: 'Table Number',
                      prefixIcon: const Icon(Icons.table_restaurant),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(16),
                      ),
                      filled: true,
                      fillColor: Theme.of(context).colorScheme.surfaceContainerHighest.withAlpha(50),
                    ),
                    keyboardType: TextInputType.number,
                  ),
                ),
              Padding(
                padding: const EdgeInsets.all(16.0),
                child: TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    hintText: 'Search items or categories...',
                    prefixIcon: const Icon(Icons.search),
                    suffixIcon: _searchQuery.isNotEmpty
                        ? IconButton(
                            icon: const Icon(Icons.clear),
                            onPressed: () {
                              _searchController.clear();
                              setState(() => _searchQuery = '');
                            },
                          )
                        : null,
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(16),
                    ),
                    filled: true,
                    fillColor: Theme.of(context).colorScheme.surfaceContainerHighest.withAlpha(50),
                  ),
                  onChanged: (value) => setState(() => _searchQuery = value),
                ),
              ),
              Expanded(
                child: filteredMenus.isEmpty
                    ? const Center(child: Text('No items found'))
                    : ListView.builder(
                        padding: const EdgeInsets.symmetric(horizontal: 16),
                        itemCount: filteredMenus.length,
                        itemBuilder: (context, index) {
                          final menu = filteredMenus[index];
                          final qty = _selectedItems[menu.id] ?? 0;

                          return Card(
                            margin: const EdgeInsets.only(bottom: 12),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(16),
                            ),
                            child: ListTile(
                              title: Text(menu.name, style: const TextStyle(fontWeight: FontWeight.bold)),
                              subtitle: Text('${menu.category} • \$${menu.price}'),
                              trailing: Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  if (qty > 0) ...[
                                    IconButton(
                                      icon: const Icon(Icons.remove_circle_outline),
                                      onPressed: () => _removeItem(menu),
                                      color: Colors.orange,
                                    ),
                                    Text('$qty', style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                                    IconButton(
                                      icon: const Icon(Icons.add_circle_outline),
                                      onPressed: () => _addItem(menu),
                                      color: Colors.orange,
                                    ),
                                  ] else
                                    IconButton(
                                      icon: const Icon(Icons.add_shopping_cart),
                                      onPressed: () => _addItem(menu),
                                      color: Theme.of(context).colorScheme.primary,
                                    ),
                                ],
                              ),
                            ),
                          );
                        },
                      ),
              ),
              if (_selectedItems.isNotEmpty)
                Container(
                  padding: const EdgeInsets.all(24),
                  decoration: BoxDecoration(
                    color: Theme.of(context).colorScheme.surfaceContainerHighest.withAlpha(100),
                    borderRadius: const BorderRadius.vertical(top: Radius.circular(32)),
                  ),
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            widget.orderId == null ? 'Create Order' : 'Add Items to Order',
                            style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                          ),
                          Text('\$${_calculateTotal(menus).toStringAsFixed(2)}',
                              style: TextStyle(
                                fontSize: 24,
                                fontWeight: FontWeight.bold,
                                color: Theme.of(context).colorScheme.primary,
                              )),
                        ],
                      ),
                      const SizedBox(height: 16),
                      SizedBox(
                        width: double.infinity,
                        child: ElevatedButton(
                          onPressed: orderAction.isLoading
                              ? null
                              : () async {
                                  final navigator = Navigator.of(context);
                                  final messenger = ScaffoldMessenger.of(context);
                                  final List<OrderItem> items = _selectedItems.entries.map<OrderItem>((entry) {
                                    final menu = menus.firstWhere((m) => m.id == entry.key);
                                    return OrderItem(
                                      id: '', // Generated by backend
                                      menuId: entry.key,
                                      menuName: menu.name,
                                      quantity: entry.value,
                                      price: menu.price.toDouble(),
                                    );
                                  }).toList();

                                  // Get the actual user ID from the auth state.
                                  final userId = ref.read(authProvider).userId ?? "00000000-0000-0000-0000-000000000000";
                                  final tableNumber = _tableNumberController.text.trim();
                                  if (tableNumber.isEmpty) {
                                    messenger.showSnackBar(
                                      const SnackBar(content: Text('Please enter table number')),
                                    );
                                    return;
                                  }

                                  if (widget.orderId == null) {
                                    await ref.read(orderActionProvider.notifier).createOrder(userId, tableNumber, items);
                                  } else {
                                    await ref.read(orderActionProvider.notifier).updateOrder(widget.orderId!, newItems: items);
                                  }
                                  
                                  if (!mounted) return;
                                  
                                  if (!ref.read(orderActionProvider).hasError) {
                                    navigator.pop();
                                    messenger.showSnackBar(
                                      SnackBar(content: Text(widget.orderId == null ? 'Order created' : 'Items added')),
                                    );
                                  }
                                },
                          style: ElevatedButton.styleFrom(
                            padding: const EdgeInsets.symmetric(vertical: 16),
                            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                          ),
                          child: orderAction.isLoading
                              ? const CircularProgressIndicator()
                              : const Text('PLACE ORDER', style: TextStyle(fontWeight: FontWeight.bold)),
                        ),
                      ),
                    ],
                  ),
                ),
            ],
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (err, stack) => Center(child: Text('Error: $err')),
      ),
    );
  }
}
