import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/user.dart';
import '../services/user_service.dart';

class UserState {
  final List<User> users;
  final bool isLoading;
  final String? error;

  UserState({
    required this.users,
    required this.isLoading,
    this.error,
  });

  factory UserState.initial() => UserState(users: [], isLoading: false);

  UserState copyWith({
    List<User>? users,
    bool? isLoading,
    String? error,
  }) {
    return UserState(
      users: users ?? this.users,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

class UserNotifier extends StateNotifier<UserState> {
  final UserService _userService;

  UserNotifier(this._userService) : super(UserState.initial()) {
    getUsers();
  }

  Future<void> getUsers() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final users = await _userService.listUsers();
      state = state.copyWith(users: users, isLoading: false);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<bool> addUser(String name, String email, String password, String role) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      await _userService.createUser(name, email, password, role);
      await getUsers();
      return true;
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
      return false;
    }
  }

  Future<bool> editUser(String id, {String? name, String? email, String? password, String? role}) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      await _userService.updateUser(id, name: name, email: email, password: password, role: role);
      await getUsers();
      return true;
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
      return false;
    }
  }

  Future<bool> removeUser(String id) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      await _userService.deleteUser(id);
      await getUsers();
      return true;
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
      return false;
    }
  }
}

final userServiceProvider = Provider((ref) => UserService());
final usersProvider = StateNotifierProvider<UserNotifier, UserState>((ref) {
  return UserNotifier(ref.watch(userServiceProvider));
});
