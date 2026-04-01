import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/auth_provider.dart';
import 'menu_page.dart';
import 'order_page.dart';
import 'user_page.dart';

class Dashboard extends ConsumerStatefulWidget {
  const Dashboard({super.key});

  @override
  ConsumerState<Dashboard> createState() => _DashboardState();
}

class _DashboardState extends ConsumerState<Dashboard> {
  int _selectedIndex = 0;

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authProvider);
    final role = authState.role ?? 'waiter';

    final List<Widget> pages = [
      const OrderPage(),
      if (role == 'admin') const MenuPage(),
      if (role == 'admin') const UserPage(),
    ];

    final List<BottomNavigationBarItem> navItems = [
      const BottomNavigationBarItem(
        icon: Icon(Icons.shopping_cart),
        label: 'Orders',
      ),
      if (role == 'admin')
        const BottomNavigationBarItem(
          icon: Icon(Icons.restaurant_menu),
          label: 'Menu',
        ),
      if (role == 'admin')
        const BottomNavigationBarItem(
          icon: Icon(Icons.people),
          label: 'Users',
        ),
    ];

    return Scaffold(
      appBar: AppBar(
        title: Text('Restaurant Pro - ${role.toUpperCase()}'),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () => ref.read(authProvider.notifier).logout(),
          ),
        ],
      ),
      body: pages[_selectedIndex],
      bottomNavigationBar: navItems.length >= 2
          ? BottomNavigationBar(
              items: navItems,
              currentIndex: _selectedIndex,
              onTap: (index) {
                setState(() {
                  _selectedIndex = index;
                });
              },
            )
          : null,
    );
  }
}
