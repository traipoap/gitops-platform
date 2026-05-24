/**
 * PDPA CRM Terminal - JavaScript functionality
 * Handles page navigation and form interactions
 */

document.addEventListener('DOMContentLoaded', function() {
    // Show a specific page by ID
    window.showPage = function(pageId) {
        // Hide all pages
        document.querySelectorAll('.page').forEach(page => {
            page.classList.remove('active');
        });
        
        // Show the specified page
        const page = document.getElementById(pageId);
        if (page) {
            page.classList.add('active');
        }
    };

    // Handle login form submission
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', function(e) {
            e.preventDefault();
            
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            
            // In a real app, this would send data to a server
            if (username && password) {
                alert('เข้าสู่ระบบสำเร็จ! การเข้าสู่ระบบสำเร็จ');
                showPage('page-dashboard');
            } else {
                alert('กรุณากรอกชื่อผู้ใช้และรหัสผ่าน');
            }
        });
    }

    // Handle forgot password button
    const recoveryEmailBtn = document.querySelector('#page-forgot .login-btn');
    if (recoveryEmailBtn) {
        recoveryEmailBtn.addEventListener('click', function(e) {
            e.preventDefault();
            
            const email = document.getElementById('recovery-email').value;
            
            if (email) {
                alert('ระบบได้ส่งอีเมลกู้คืนรหัสผ่านไปยังอีเมลของคุณ\nPassword recovery email has been sent');
                showPage('page-login');
            } else {
                alert('กรุณากรอกอีเมลของคุณ');
            }
        });
    }

    // Handle register button
    const registerBtn = document.querySelector('#page-register .login-btn');
    if (registerBtn) {
        registerBtn.addEventListener('click', function(e) {
            e.preventDefault();
            
            const password = document.getElementById('reg-password').value;
            const confirmPassword = document.getElementById('reg-confirm').value;
            
            if (password && confirmPassword && password === confirmPassword) {
                alert('ลงทะเบียนสำเร็จ! ระบบจะส่งอีเมลยืนยันไปยังอีเมลของคุณ\nRegistration successful! Verification email sent.');
                showPage('page-login');
            } else {
                alert('รหัสผ่านไม่ตรงกัน กรุณาตรวจสอบ');
            }
        });
    }
});