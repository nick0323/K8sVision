#!/usr/bin/env node

/**
 * 批量重构页面组件脚本
 * 将所有页面组件重构为使用ResourcePage通用组件
 */

const fs = require('fs');
const path = require('path');

// 页面配置映射
const PAGE_CONFIGS = {
  'StatefulSetsPage.jsx': 'STATEFULSETS_CONFIG',
  'DaemonSetsPage.jsx': 'DAEMONSETS_CONFIG',
  'CronJobsPage.jsx': 'CRONJOBS_CONFIG',
  'JobsPage.jsx': 'JOBS_CONFIG',
  'IngressPage.jsx': 'INGRESS_CONFIG',
  'ServicesPage.jsx': 'SERVICES_CONFIG',
  'EventsPage.jsx': 'EVENTS_CONFIG',
  'PVCsPage.jsx': 'PVCS_CONFIG',
  'PVsPage.jsx': 'PVS_CONFIG',
  'StorageClassesPage.jsx': 'STORAGECLASSES_CONFIG',
  'ConfigMapsPage.jsx': 'CONFIGMAPS_CONFIG',
  'SecretsPage.jsx': 'SECRETS_CONFIG',
  'NamespacesPage.jsx': 'NAMESPACES_CONFIG',
  'NodesPage.jsx': 'NODES_CONFIG'
};

// 新组件模板
const NEW_COMPONENT_TEMPLATE = (configName) => `import React from 'react';
import ResourcePage from './components/ResourcePage';
import { ${configName} } from './constants/pageConfigs';

export default function {PAGE_NAME}({ collapsed, onToggleCollapsed }) {
  return (
    <ResourcePage
      {...${configName}}
      collapsed={collapsed}
      onToggleCollapsed={onToggleCollapsed}
    />
  );
}`;

// 获取页面名称
function getPageName(filename) {
  return filename.replace('Page.jsx', '');
}

// 重构单个页面
function refactorPage(filePath, configName) {
  const pageName = getPageName(path.basename(filePath));
  const newContent = NEW_COMPONENT_TEMPLATE(configName).replace('{PAGE_NAME}', pageName);
  
  try {
    fs.writeFileSync(filePath, newContent, 'utf8');
    console.log(`✅ 重构完成: ${filePath}`);
  } catch (error) {
    console.error(`❌ 重构失败: ${filePath}`, error.message);
  }
}

// 主函数
function main() {
  const srcDir = path.join(__dirname, '..', 'src');
  
  console.log('🚀 开始批量重构页面组件...\n');
  
  Object.entries(PAGE_CONFIGS).forEach(([filename, configName]) => {
    const filePath = path.join(srcDir, filename);
    
    if (fs.existsSync(filePath)) {
      console.log(`🔄 重构: ${filename}`);
      refactorPage(filePath, configName);
    } else {
      console.log(`⚠️  文件不存在: ${filename}`);
    }
  });
  
  console.log('\n🎉 批量重构完成!');
  console.log('\n📝 注意事项:');
  console.log('1. 确保constants/pageConfigs.js文件已创建');
  console.log('2. 确保components/ResourcePage.jsx文件已创建');
  console.log('3. 检查所有页面是否正确导入');
}

if (require.main === module) {
  main();
}

module.exports = { refactorPage, PAGE_CONFIGS };

