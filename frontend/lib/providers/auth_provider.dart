import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../services/auth_service.dart';

class AuthState {
  final bool isAuthenticated;
  final String? role;
  final String? userId;
  final String? error;
  final bool isLoading;

  AuthState({
    required this.isAuthenticated,
    this.role,
    this.userId,
    this.error,
    required this.isLoading,
  });

  factory AuthState.initial() => AuthState(isAuthenticated: false, isLoading: false);

  AuthState copyWith({
    bool? isAuthenticated,
    String? role,
    String? userId,
    String? error,
    bool? isLoading,
  }) {
    return AuthState(
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      role: role ?? this.role,
      userId: userId ?? this.userId,
      error: error ?? this.error,
      isLoading: isLoading ?? this.isLoading,
    );
  }
}

class AuthNotifier extends StateNotifier<AuthState> {
  final AuthService _authService;

  AuthNotifier(this._authService) : super(AuthState.initial()) {
    checkStatus();
  }

  Future<void> checkStatus() async {
    state = state.copyWith(isLoading: true);
    final isLoggedIn = await _authService.isLoggedIn();
    if (isLoggedIn) {
      final role = await _authService.getRole();
      final userId = await _authService.getUserId();
      state = state.copyWith(isAuthenticated: true, role: role, userId: userId, isLoading: false);
    } else {
      state = state.copyWith(isAuthenticated: false, isLoading: false);
    }
  }

  Future<bool> login(String email, String password) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final response = await _authService.login(email, password);
      if (response.success) {
        state = state.copyWith(
          isAuthenticated: true,
          role: response.role,
          userId: response.userId,
          isLoading: false,
        );
        return true;
      } else {
        state = state.copyWith(error: response.errMsg, isLoading: false);
        return false;
      }
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
      return false;
    }
  }

  Future<void> logout() async {
    await _authService.logout();
    state = AuthState.initial();
  }
}

final authServiceProvider = Provider((ref) => AuthService());
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier(ref.watch(authServiceProvider));
});
