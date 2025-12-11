<template>
  <div class="container">
    <div class="header">
      <h1>🚀 Gomoco</h1>
      <p>轻量级 Mock Server 管理平台</p>
    </div>

    <!-- Alert Messages -->
    <div v-if="alert.show" :class="['alert', alert.type === 'success' ? 'alert-success' : 'alert-error']">
      {{ alert.message }}
    </div>

    <!-- Create/Edit Form -->
    <div class="form-section">
      <h2>{{ editingMock ? '编辑 Mock API' : '创建 Mock API' }}</h2>
      <form @submit.prevent="submitForm">
        <div class="form-group">
          <label for="name">API 名称 *</label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            required
            placeholder="例如: 用户登录接口"
          />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="port">端口 *</label>
            <input
              id="port"
              v-model.number="form.port"
              type="number"
              min="1"
              max="65535"
              required
              :disabled="!!editingMock"
              placeholder="例如: 9090"
            />
          </div>

          <div class="form-group">
            <label for="protocol">协议 *</label>
            <select id="protocol" v-model="form.protocol" required :disabled="!!editingMock">
              <option value="http">HTTP</option>
              <option value="tcp">TCP</option>
            </select>
          </div>

          <div class="form-group">
            <label for="charset">字符集 *</label>
            <select id="charset" v-model="form.charset" required>
              <option value="UTF-8">UTF-8</option>
              <option value="GBK">GBK</option>
            </select>
          </div>
        </div>

        <div v-if="form.protocol === 'http'" class="form-row">
          <div class="form-group">
            <label for="path">路径</label>
            <input
              id="path"
              v-model="form.path"
              type="text"
              placeholder="例如: /api/test (默认为 /)"
            />
          </div>

          <div class="form-group">
            <label for="method">HTTP 方法</label>
            <select id="method" v-model="form.method">
              <option value="">任意方法</option>
              <option value="GET">GET</option>
              <option value="POST">POST</option>
              <option value="PUT">PUT</option>
              <option value="DELETE">DELETE</option>
            </select>
          </div>
        </div>

        <div class="form-group">
          <label for="content">响应内容 *</label>
          <textarea
            id="content"
            v-model="form.content"
            required
            placeholder="输入固定返回的报文内容..."
          ></textarea>
        </div>

        <div style="display: flex; gap: 10px;">
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ editingMock ? '更新' : '创建' }} Mock API
          </button>
          <button
            v-if="editingMock"
            type="button"
            class="btn btn-warning"
            @click="cancelEdit"
          >
            取消编辑
          </button>
        </div>
      </form>
    </div>

    <!-- Mock List -->
    <div class="list-section">
      <h2>Mock API 列表</h2>
      
      <div v-if="loading && mocks.length === 0" class="loading">
        加载中...
      </div>

      <div v-else-if="mocks.length === 0" class="empty-state">
        <p style="font-size: 1.2rem;">暂无 Mock API</p>
        <p>创建第一个 Mock API 开始使用</p>
      </div>

      <div v-else class="mock-list">
        <div v-for="mock in mocks" :key="mock.id" class="mock-item">
          <div class="mock-header">
            <div class="mock-title">
              {{ mock.name }}
            </div>
            <span :class="['status-badge', mock.status === 'running' ? 'status-running' : 'status-stopped']">
              {{ mock.status === 'running' ? '运行中' : '已停止' }}
            </span>
          </div>

          <div class="mock-details">
            <div class="detail-item">
              <span class="detail-label">名称</span>
              <span class="detail-value">{{ mock.name }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">协议</span>
              <span class="detail-value">{{ mock.protocol.toUpperCase() }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">端口</span>
              <span class="detail-value">{{ mock.port }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">字符集</span>
              <span class="detail-value">{{ mock.charset }}</span>
            </div>
            <div v-if="mock.protocol === 'http' && mock.path" class="detail-item">
              <span class="detail-label">路径</span>
              <span class="detail-value">{{ mock.path }}</span>
            </div>
            <div v-if="mock.protocol === 'http' && mock.method" class="detail-item">
              <span class="detail-label">方法</span>
              <span class="detail-value">{{ mock.method }}</span>
            </div>
          </div>

          <div class="mock-content">{{ mock.content }}</div>

          <div class="mock-actions">
            <button class="btn btn-success" @click="editMock(mock)">编辑</button>
            <button class="btn btn-danger" @click="deleteMock(mock.id)">删除</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'App',
  data() {
    return {
      mocks: [],
      loading: false,
      editingMock: null,
      form: {
        name: '',
        port: '',
        protocol: 'http',
        content: '',
        charset: 'UTF-8',
        path: '',
        method: ''
      },
      alert: {
        show: false,
        type: 'success',
        message: ''
      }
    }
  },
  mounted() {
    this.loadMocks()
  },
  methods: {
    async loadMocks() {
      try {
        this.loading = true
        const response = await axios.get('/api/mocks')
        this.mocks = response.data || []
      } catch (error) {
        this.showAlert('error', '加载 Mock API 列表失败: ' + error.message)
      } finally {
        this.loading = false
      }
    },
    async submitForm() {
      try {
        this.loading = true
        
        if (this.editingMock) {
          // Update existing mock
          await axios.put(`/api/mocks/${this.editingMock.id}`, {
            name: this.form.name,
            content: this.form.content,
            charset: this.form.charset,
            path: this.form.path,
            method: this.form.method
          })
          this.showAlert('success', 'Mock API 更新成功!')
        } else {
          // Create new mock
          await axios.post('/api/mocks', this.form)
          this.showAlert('success', 'Mock API 创建成功!')
        }
        
        this.resetForm()
        await this.loadMocks()
      } catch (error) {
        this.showAlert('error', '操作失败: ' + (error.response?.data?.error || error.message))
      } finally {
        this.loading = false
      }
    },
    editMock(mock) {
      this.editingMock = mock
      this.form = {
        name: mock.name,
        port: mock.port,
        protocol: mock.protocol,
        content: mock.content,
        charset: mock.charset,
        path: mock.path || '',
        method: mock.method || ''
      }
      window.scrollTo({ top: 0, behavior: 'smooth' })
    },
    cancelEdit() {
      this.resetForm()
    },
    async deleteMock(id) {
      if (!confirm('确定要删除这个 Mock API 吗?')) {
        return
      }
      
      try {
        this.loading = true
        await axios.delete(`/api/mocks/${id}`)
        this.showAlert('success', 'Mock API 删除成功!')
        await this.loadMocks()
      } catch (error) {
        this.showAlert('error', '删除失败: ' + (error.response?.data?.error || error.message))
      } finally {
        this.loading = false
      }
    },
    resetForm() {
      this.editingMock = null
      this.form = {
        name: '',
        port: '',
        protocol: 'http',
        content: '',
        charset: 'UTF-8',
        path: '',
        method: ''
      }
    },
    showAlert(type, message) {
      this.alert = { show: true, type, message }
      setTimeout(() => {
        this.alert.show = false
      }, 5000)
    }
  }
}
</script>
