// 1. 导入封装好的 service 实例
import request from './request'

/**
 * 获取文章列表
 * @param {object} params - 查询参数，例如 { page: 1, pageSize: 10 }
 */
export function fetchPostList(params) {
    return request({
        url: '/posts',
        method: 'get',
        params: params, // get 请求的查询参数要放在 params 里
    })
}

/**
 * 根据 ID 获取文章详情
 * @param {number} id - 文章的 ID
 */
export function fetchPostById(id) {
    return request({
        url: `/posts/${id}`, // 使用模板字符串拼接 URL
        method: 'get',
    })
}