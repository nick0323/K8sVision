#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

// 定义CSS变量映射
const CSS_VARIABLES = {
  // 颜色变量
  '--primary': '#1890ff',
  '--primary-light': '#e6f7ff',
  '--success': '#52c41a',
  '--warning': '#faad14',
  '--danger': '#ff4d4f',
  '--text-primary': '#222',
  '--text-secondary': '#444',
  '--text-tertiary': '#666',
  '--text-quaternary': '#888',
  '--text-disabled': '#bbb',
  '--bg-primary': '#f6f8fa',
  '--bg-secondary': '#fff',
  '--bg-tertiary': '#fafbfc',
  '--bg-container': '#f0f5ff',
  '--bg-tag': '#f8fafc',
  '--border-primary': '#e5e7eb',
  '--border-secondary': '#f0f0f0',
  '--border-container': '#d6e4ff',
  '--border-tag': '#e2e8f0',
  
  // 字体大小变量
  '--font-size-xs': '11px',
  '--font-size-sm': '12px',
  '--font-size-base': '14px',
  '--font-size-lg': '15px',
  '--font-size-xl': '16px',
  '--font-size-2xl': '18px',
  '--font-size-3xl': '22px',
  '--font-size-4xl': '26px',
  '--font-size-5xl': '32px',
  '--font-size-6xl': '48px',
  
  // 间距变量
  '--spacing-xs': '2px',
  '--spacing-sm': '8px',
  '--spacing-md': '16px',
  '--spacing-lg': '24px',
  '--spacing-xl': '32px',
  
  // 圆角变量
  '--radius-xs': '4px',
  '--radius-sm': '6px',
  '--radius-md': '10px',
  '--radius-lg': '14px',
  
  // 字体变量
  '--font-main': "'Maple Mono', monospace"
};

// 硬编码值映射到建议的变量
const HARDCODED_SUGGESTIONS = {
  '#1890ff': 'var(--primary)',
  '#e6f7ff': 'var(--primary-light)',
  '#52c41a': 'var(--success)',
  '#faad14': 'var(--warning)',
  '#ff4d4f': 'var(--danger)',
  '#222': 'var(--text-primary)',
  '#444': 'var(--text-secondary)',
  '#666': 'var(--text-tertiary)',
  '#888': 'var(--text-quaternary)',
  '#bbb': 'var(--text-disabled)',
  '#f6f8fa': 'var(--bg-primary)',
  '#fff': 'var(--bg-secondary)',
  '#fafbfc': 'var(--bg-tertiary)',
  '#f0f5ff': 'var(--bg-container)',
  '#f8fafc': 'var(--bg-tag)',
  '#e5e7eb': 'var(--border-primary)',
  '#f0f0f0': 'var(--border-secondary)',
  '#d6e4ff': 'var(--border-container)',
  '#e2e8f0': 'var(--border-tag)',
  '11px': 'var(--font-size-xs)',
  '12px': 'var(--font-size-sm)',
  '14px': 'var(--font-size-base)',
  '15px': 'var(--font-size-lg)',
  '16px': 'var(--font-size-xl)',
  '18px': 'var(--font-size-2xl)',
  '22px': 'var(--font-size-3xl)',
  '26px': 'var(--font-size-4xl)',
  '32px': 'var(--font-size-5xl)',
  '48px': 'var(--font-size-6xl)',
  '2px': 'var(--spacing-xs)',
  '8px': 'var(--spacing-sm)',
  '16px': 'var(--spacing-md)',
  '24px': 'var(--spacing-lg)',
  '32px': 'var(--spacing-xl)',
  '4px': 'var(--radius-xs)',
  '6px': 'var(--radius-sm)',
  '10px': 'var(--radius-md)',
  '14px': 'var(--radius-lg)',
  "'Maple Mono', monospace": 'var(--font-main)'
};

function checkCSSFile(filePath) {
  const content = fs.readFileSync(filePath, 'utf8');
  const lines = content.split('\n');
  const issues = [];
  
  lines.forEach((line, index) => {
    // 跳过注释和变量定义
    if (line.trim().startsWith('/*') || line.trim().startsWith('*') || line.includes('--')) {
      return;
    }
    
    // 检查硬编码的颜色值
    Object.keys(HARDCODED_SUGGESTIONS).forEach(hardcoded => {
      if (line.includes(hardcoded)) {
        issues.push({
          line: index + 1,
          value: hardcoded,
          suggestion: HARDCODED_SUGGESTIONS[hardcoded],
          context: line.trim()
        });
      }
    });
  });
  
  return issues;
}

function scanDirectory(dir) {
  const files = fs.readdirSync(dir);
  const cssFiles = [];
  
  files.forEach(file => {
    const filePath = path.join(dir, file);
    const stat = fs.statSync(filePath);
    
    if (stat.isDirectory() && !file.startsWith('.') && file !== 'node_modules') {
      cssFiles.push(...scanDirectory(filePath));
    } else if (file.endsWith('.css')) {
      cssFiles.push(filePath);
    }
  });
  
  return cssFiles;
}

function main() {
  const srcDir = path.join(__dirname, '..', 'frontend', 'src');
  const cssFiles = scanDirectory(srcDir);
  
  console.log('🔍 检查CSS变量使用情况...\n');
  
  let totalIssues = 0;
  
  cssFiles.forEach(file => {
    const issues = checkCSSFile(file);
    if (issues.length > 0) {
      console.log(`📁 ${path.relative(process.cwd(), file)}`);
      issues.forEach(issue => {
        console.log(`  Line ${issue.line}: ${issue.value} → ${issue.suggestion}`);
        console.log(`    Context: ${issue.context}`);
      });
      console.log('');
      totalIssues += issues.length;
    }
  });
  
  if (totalIssues === 0) {
    console.log('✅ 所有CSS文件都正确使用了变量！');
  } else {
    console.log(`⚠️  发现 ${totalIssues} 个建议优化的地方`);
    console.log('💡 建议使用CSS变量替代硬编码值以提高一致性和可维护性');
  }
}

if (require.main === module) {
  main();
}

module.exports = { checkCSSFile, scanDirectory, CSS_VARIABLES, HARDCODED_SUGGESTIONS }; 