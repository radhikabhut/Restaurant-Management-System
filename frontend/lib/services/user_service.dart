import 'dart:convert';
import '../models/user.dart';
import 'api_client.dart';

class UserService {
  final ApiClient _apiClient = ApiClient();

  Future<List<User>> listUsers() async {
    final response = await _apiClient.post('/users/list', {'limit': 100, 'offset': 0});
    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final List<dynamic> usersJson = data['users'] ?? [];
      return usersJson.map((json) => User.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load users: ${response.statusCode}');
    }
  }

  Future<void> createUser(String name, String email, String password, String role) async {
    final response = await _apiClient.post('/users/create', {
      'name': name,
      'email': email,
      'password': password,
      'role': role,
    });
    if (response.statusCode != 200) {
      final data = jsonDecode(response.body);
      throw Exception(data['errMsg'] ?? 'Failed to create user');
    }
  }

  Future<void> updateUser(String id, {String? name, String? email, String? password, String? role}) async {
    final Map<String, dynamic> body = {
      'id': id,
    };
    if (name != null) body['name'] = name;
    if (email != null) body['email'] = email;
    if (password != null) body['password'] = password;
    if (role != null) body['role'] = role;

    final response = await _apiClient.post('/users/update', body);
    if (response.statusCode != 200) {
      final data = jsonDecode(response.body);
      throw Exception(data['errMsg'] ?? 'Failed to update user');
    }
  }

  Future<void> deleteUser(String id) async {
    final response = await _apiClient.post('/users/delete', {'id': id});
    if (response.statusCode != 200) {
      final data = jsonDecode(response.body);
      throw Exception(data['errMsg'] ?? 'Failed to delete user');
    }
  }
}
