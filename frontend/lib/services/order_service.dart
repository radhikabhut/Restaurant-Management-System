import 'dart:convert';
import '../models/order.dart';
import 'api_client.dart';

class OrderService {
  final ApiClient _apiClient = ApiClient();

  Future<List<Order>> listOrders({String? userId, String? status}) async {
    final Map<String, dynamic> body = {};
    if (userId != null) body['userId'] = userId;
    if (status != null) body['status'] = status;

    final response = await _apiClient.post('/orders/list', body);

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final List<dynamic> ordersJson = data['orders'] ?? [];
      return ordersJson.map((json) => Order.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load orders');
    }
  }

  Future<String> createOrder(String userId, String tableNumber, List<OrderItem> items) async {
    final response = await _apiClient.post('/orders/create', {
      'userId': userId,
      'tableNumber': tableNumber,
      'orderItems': items.map((it) => {
        'menuId': it.menuId,
        'quantity': it.quantity,
      }).toList(),
    });

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return data['id'];
    } else {
      throw Exception('Failed to create order');
    }
  }

  Future<void> updateOrder(String id, {String? status, List<OrderItem>? newItems, List<String>? deletedItemIds}) async {
    final Map<String, dynamic> body = {'id': id};
    if (status != null) body['status'] = status;
    if (newItems != null && newItems.isNotEmpty) {
      body['orderItems'] = newItems.map((it) => {
        'menuId': it.menuId,
        'quantity': it.quantity,
      }).toList();
    }
    if (deletedItemIds != null && deletedItemIds.isNotEmpty) {
      body['deletedItemIds'] = deletedItemIds;
    }

    final response = await _apiClient.post('/orders/update', body);

    if (response.statusCode != 200) {
      throw Exception('Failed to update order');
    }
  }

  Future<Map<String, dynamic>> generateBill(String orderId) async {
    final response = await _apiClient.post('/orders/bill', {
      'orderId': orderId,
    });

    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    } else {
      throw Exception('Failed to generate bill');
    }
  }

  Future<void> deleteOrder(String id, String userId) async {
    final response = await _apiClient.post('/orders/delete', {
      'id': id,
      'userId': userId,
    });

    if (response.statusCode != 200) {
      throw Exception('Failed to delete order');
    }
  }
}
