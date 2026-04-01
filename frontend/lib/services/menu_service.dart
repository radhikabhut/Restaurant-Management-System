import 'dart:convert';
import '../models/menu_item.dart';
import 'api_client.dart';

class MenuService {
  final ApiClient _apiClient = ApiClient();

  Future<List<MenuItem>> listMenus({int limit = 10, int offset = 0}) async {
    final response = await _apiClient.post('/menus/list', {
      'limit': limit,
      'offset': offset,
    });

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final List<dynamic> menusJson = data['menus'] ?? [];
      return menusJson.map((json) => MenuItem.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load menus');
    }
  }

  Future<void> createMenu(MenuItem menu) async {
    final response = await _apiClient.post('/menus/create', {
      'name': menu.name,
      'price': menu.price,
      'category': menu.category,
      'isAvailable': menu.isAvailable,
    });

    if (response.statusCode != 200) {
      throw Exception('Failed to create menu');
    }
  }

  Future<void> updateMenu(MenuItem menu) async {
    final response = await _apiClient.post('/menus/update', {
      'id': menu.id,
      'name': menu.name,
      'price': menu.price,
      'category': menu.category,
      'isAvailable': menu.isAvailable,
    });

    if (response.statusCode != 200) {
      throw Exception('Failed to update menu');
    }
  }

  Future<void> deleteMenu(String id) async {
    final response = await _apiClient.post('/menus/delete', {'id': id});

    if (response.statusCode != 200) {
      throw Exception('Failed to delete menu');
    }
  }
}
