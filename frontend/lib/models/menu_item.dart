import 'package:json_annotation/json_annotation.dart';

part 'menu_item.g.dart';

@JsonSerializable()
class MenuItem {
  final String id;
  final String name;
  final int price;
  final String category;
  final bool isAvailable;
  final String? createdAt;
  final String? updatedAt;

  MenuItem({
    required this.id,
    required this.name,
    required this.price,
    required this.category,
    required this.isAvailable,
    this.createdAt,
    this.updatedAt,
  });

  factory MenuItem.fromJson(Map<String, dynamic> json) => _$MenuItemFromJson(json);
  Map<String, dynamic> toJson() => _$MenuItemToJson(this);
}
