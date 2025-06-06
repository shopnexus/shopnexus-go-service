đổi custom_permission cho staff và admin only
đổi lại đối tượng người dùng: người dùng và các quản trị viên
xoá bảng Staff, chỉ còn User/Admin
xài account.type thay vì account.role để phân loại User/Admin
sửa account role thành có thể có nhiều roles thay vì 1
di chuyển avatar_url sang User + Admin thay vì để ở base "Account"
xoá gần hết date_updated
product đc chia ra thêm 1 table ProductTracking, chỉ chứa current_stock và sold; product.add_price đổi tên thành additional_price; product.quantity trước đây mình kiểm tra bằng cách quantity - sold > 0; giờ phải sửa lại thành chỉ cần current_stock > 0
sale xoá hết field id ràng buộc và dùng item_id để polymorphsm theo SaleType enum, sale cũng chia thêm 1 table SaleTracking gồm current_stock và used
sửa lại toàn bộ field trong PaymentVnpay
refund thêm trường "amount" để biết số tiền khách sẽ đc nhận, thay vì phụ thuộc hoàn toàn vào productonpayment.total_price, vì admin có thể thêm penalty fee vào refund
