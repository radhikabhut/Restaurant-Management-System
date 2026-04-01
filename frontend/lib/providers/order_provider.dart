import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/order.dart';
import '../services/order_service.dart';

final orderServiceProvider = Provider((ref) => OrderService());

final orderListProvider = FutureProvider.family<List<Order>, ({String? userId, String? status})>((ref, arg) async {
  final service = ref.watch(orderServiceProvider);
  return service.listOrders(userId: arg.userId, status: arg.status);
});

class OrderNotifier extends StateNotifier<AsyncValue<void>> {
  final OrderService _service;
  final Ref _ref;

  OrderNotifier(this._service, this._ref) : super(const AsyncValue.data(null));

  Future<void> createOrder(String userId, String tableNumber, List<OrderItem> items) async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() async {
      await _service.createOrder(userId, tableNumber, items);
      _ref.invalidate(orderListProvider);
    });
  }

  Future<void> updateOrder(String id, {String? status, List<OrderItem>? newItems, List<String>? deletedItemIds}) async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() async {
      await _service.updateOrder(id, status: status, newItems: newItems, deletedItemIds: deletedItemIds);
      _ref.invalidate(orderListProvider);
    });
  }

  Future<void> deleteOrder(String id, String userId) async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() async {
      await _service.deleteOrder(id, userId);
      _ref.invalidate(orderListProvider);
    });
  }

  Future<Map<String, dynamic>> generateBill(String orderId) async {
    return await _service.generateBill(orderId);
  }
}

final orderActionProvider = StateNotifierProvider<OrderNotifier, AsyncValue<void>>((ref) {
  return OrderNotifier(ref.watch(orderServiceProvider), ref);
});
