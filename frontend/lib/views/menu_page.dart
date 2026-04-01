import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/menu_provider.dart';
import '../models/menu_item.dart';
import '../providers/auth_provider.dart';

class MenuPage extends ConsumerWidget {
  const MenuPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final menuList = ref.watch(menuListProvider);
    final authState = ref.watch(authProvider);
    final isAdmin = authState.role == 'admin';

    return Scaffold(
      body: menuList.when(
        data: (menus) => ListView.builder(
          padding: const EdgeInsets.all(16),
          itemCount: menus.length,
          itemBuilder: (context, index) {
            final menu = menus[index];
            return Card(
              margin: const EdgeInsets.only(bottom: 12),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(16),
              ),
              child: ListTile(
                contentPadding: const EdgeInsets.all(16),
                title: Text(
                  menu.name,
                  style: const TextStyle(fontWeight: FontWeight.bold),
                ),
                subtitle: Text('${menu.category} • \$${menu.price}'),
                trailing: Chip(
                  label: Text(menu.isAvailable ? 'Available' : 'Sold Out'),
                  backgroundColor: menu.isAvailable
                      ? Colors.green.withAlpha(50)
                      : Colors.red.withAlpha(50),
                ),
                onTap: isAdmin ? () => _showEditDialog(context, ref, menu) : null,
              ),
            );
          },
        ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (err, stack) => Center(child: Text('Error: $err')),
      ),
      floatingActionButton: isAdmin
          ? FloatingActionButton(
              onPressed: () => _showAddDialog(context, ref),
              child: const Icon(Icons.add),
            )
          : null,
    );
  }

  void _showAddDialog(BuildContext context, WidgetRef ref) {
    _showMenuDialog(context, ref);
  }

  void _showEditDialog(BuildContext context, WidgetRef ref, MenuItem menu) {
    _showMenuDialog(context, ref, menu: menu);
  }

  void _showMenuDialog(BuildContext context, WidgetRef ref, {MenuItem? menu}) {
    final nameController = TextEditingController(text: menu?.name);
    final priceController = TextEditingController(text: menu?.price.toString());
    final categoryController = TextEditingController(text: menu?.category);
    bool isAvailable = menu?.isAvailable ?? true;
    final formKey = GlobalKey<FormState>();

    showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          title: Text(menu == null ? 'Add Menu Item' : 'Edit Menu Item'),
          content: Form(
            key: formKey,
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  TextFormField(
                    controller: nameController,
                    decoration: const InputDecoration(labelText: 'Name'),
                    validator: (v) => v == null || v.isEmpty ? 'Required' : null,
                  ),
                  TextFormField(
                    controller: priceController,
                    decoration: const InputDecoration(labelText: 'Price'),
                    keyboardType: TextInputType.number,
                    validator: (v) => v == null || double.tryParse(v) == null ? 'Invalid price' : null,
                  ),
                  TextFormField(
                    controller: categoryController,
                    decoration: const InputDecoration(labelText: 'Category'),
                    validator: (v) => v == null || v.isEmpty ? 'Required' : null,
                  ),
                  SwitchListTile(
                    title: const Text('Available'),
                    value: isAvailable,
                    onChanged: (v) => setDialogState(() => isAvailable = v),
                    contentPadding: EdgeInsets.zero,
                  ),
                ],
              ),
            ),
          ),
          actions: [
            if (menu != null)
              TextButton(
                onPressed: () => _showDeleteConfirmation(context, ref, menu),
                child: const Text('Delete', style: TextStyle(color: Colors.red)),
              ),
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () async {
                if (formKey.currentState!.validate()) {
                  final newItem = MenuItem(
                    id: menu?.id ?? '',
                    name: nameController.text,
                    price: int.parse(priceController.text), // Backend expects int64 (cents or unit)
                    category: categoryController.text,
                    isAvailable: isAvailable,
                  );

                  if (menu == null) {
                    await ref.read(menuActionProvider.notifier).createMenu(newItem);
                  } else {
                    await ref.read(menuActionProvider.notifier).updateMenu(newItem);
                  }

                  if (context.mounted) Navigator.pop(context);
                }
              },
              child: Text(menu == null ? 'Create' : 'Save'),
            ),
          ],
        ),
      ),
    );
  }

  void _showDeleteConfirmation(BuildContext context, WidgetRef ref, MenuItem menu) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Item'),
        content: Text('Are you sure you want to delete ${menu.name}?'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context), child: const Text('Cancel')),
          TextButton(
            onPressed: () async {
              await ref.read(menuActionProvider.notifier).deleteMenu(menu.id);
              if (context.mounted) {
                Navigator.pop(context); // Close confirmation
                Navigator.pop(context); // Close edit dialog
              }
            },
            child: const Text('Delete', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  }
}
