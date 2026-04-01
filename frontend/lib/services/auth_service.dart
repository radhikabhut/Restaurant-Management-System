import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/auth.dart';
import 'api_client.dart';

class AuthService {
  final ApiClient _apiClient = ApiClient();

  Future<LoginResponse> login(String email, String password) async {
    final request = LoginRequest(email: email, password: password);
    final response = await _apiClient.post('/login', request.toJson());

    if (response.statusCode == 200) {
      final loginRes = LoginResponse.fromJson(jsonDecode(response.body));
      if (loginRes.success) {
        final prefs = await SharedPreferences.getInstance();
        await prefs.setString('auth_token', loginRes.token);
        await prefs.setString('user_role', loginRes.role);
        await prefs.setString('user_id', loginRes.userId);
      }
      return loginRes;
    } else {
      throw Exception('Login failed: ${response.statusCode}');
    }
  }

  Future<void> logout() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('auth_token');
    await prefs.remove('user_role');
    await prefs.remove('user_id');
  }

  Future<String?> getRole() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('user_role');
  }

  Future<String?> getUserId() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('user_id');
  }

  Future<bool> isLoggedIn() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.containsKey('auth_token');
  }
}
