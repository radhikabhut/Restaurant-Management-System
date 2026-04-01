import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/menu_item.dart';
import '../services/menu_service.dart';

final menuServiceProvider = Provider((ref) => MenuService());

final menuListProvider = FutureProvider<List<MenuItem>>((ref) async {
  final service = ref.watch(menuServiceProvider);
  return service.listMenus();
});

class MenuNotifier extends StateNotifier<AsyncValue<void>> {
  final MenuService _service;
  final Ref _ref;

  MenuNotifier(this._service, this._ref) : super(const AsyncValue.data(null));

  Future<void> createMenu(MenuItem menu) async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() async {
      await _service.createMenu(menu);
      _ref.invalidate(menuListProvider);
    });
  }

  Future<void> updateMenu(MenuItem menu) async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() async {
      await _service.updateMenu(menu);
      _ref.invalidate(menuListProvider);
    });
  }

  Future<void> deleteMenu(String id) async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() async {
      await _service.deleteMenu(id);
      _ref.invalidate(menuListProvider);
    });
  }
}

final menuActionProvider = StateNotifierProvider<MenuNotifier, AsyncValue<void>>((ref) {
  return MenuNotifier(ref.watch(menuServiceProvider), ref);
});
