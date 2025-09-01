// 导入封装好的 axios 实例
import request from './request'

/**
 * @description 发送用户登录请求
 * @param {object} data - 包含用户名和密码的对象, e.g., { username: 'admin', password: 'password123' }
 * @returns {Promise} - 返回一个 Promise 对象，成功时包含 token 等数据
 */
export function login(data) {
    return request({
        url: '/login', // 后端登录接口的路径
        method: 'post',
        data: data, // post 请求的数据放在 data 中
    })
}

/**
 * @description 发送用户注册请求
 * @param {object} data - 包含用户名和密码的对象, e.g., { username: 'newuser', password: 'newpassword' }
 * @returns {Promise}
 */
export function register(data) {
    return request({
        url: 'signup', 
        method: 'post',
        data: data,
    })
}
