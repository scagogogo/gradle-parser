#!/usr/bin/env node

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Get all markdown files
function getAllMarkdownFiles(dir, files = []) {
  const items = fs.readdirSync(dir);
  
  for (const item of items) {
    const fullPath = path.join(dir, item);
    const stat = fs.statSync(fullPath);
    
    if (stat.isDirectory() && !item.startsWith('.') && item !== 'node_modules') {
      getAllMarkdownFiles(fullPath, files);
    } else if (item.endsWith('.md')) {
      files.push(fullPath);
    }
  }
  
  return files;
}

// Extract links from markdown content
function extractLinks(content) {
  const linkRegex = /\[([^\]]*)\]\(([^)]+)\)/g;
  const links = [];
  let match;
  
  while ((match = linkRegex.exec(content)) !== null) {
    const [, text, url] = match;
    if (!url.startsWith('http') && !url.startsWith('mailto:')) {
      links.push({ text, url, line: content.substring(0, match.index).split('\n').length });
    }
  }
  
  return links;
}

// Check if a file exists
function checkFileExists(filePath, basePath) {
  // Remove anchor
  const cleanPath = filePath.split('#')[0];
  
  // Handle relative paths
  let fullPath;
  if (cleanPath.startsWith('./')) {
    fullPath = path.resolve(basePath, cleanPath);
  } else if (cleanPath.startsWith('../')) {
    fullPath = path.resolve(basePath, cleanPath);
  } else if (cleanPath.startsWith('/')) {
    fullPath = path.resolve(path.join(__dirname, '..'), cleanPath.substring(1));
  } else {
    fullPath = path.resolve(basePath, cleanPath);
  }
  
  // Add .md extension if not present and not a directory
  if (!path.extname(fullPath) && !fs.existsSync(fullPath)) {
    fullPath += '.md';
  }
  
  return fs.existsSync(fullPath);
}

// Main function
function checkLinks() {
  const docsDir = path.join(__dirname, '..');
  const markdownFiles = getAllMarkdownFiles(docsDir);
  
  let totalLinks = 0;
  let brokenLinks = 0;
  const issues = [];
  
  console.log('🔍 Checking internal links in markdown files...\n');
  
  for (const file of markdownFiles) {
    const content = fs.readFileSync(file, 'utf-8');
    const links = extractLinks(content);
    const basePath = path.dirname(file);
    
    totalLinks += links.length;
    
    for (const link of links) {
      if (!checkFileExists(link.url, basePath)) {
        brokenLinks++;
        const relativePath = path.relative(docsDir, file);
        issues.push({
          file: relativePath,
          line: link.line,
          text: link.text,
          url: link.url
        });
      }
    }
  }
  
  // Report results
  console.log(`📊 Link Check Results:`);
  console.log(`   Total files: ${markdownFiles.length}`);
  console.log(`   Total links: ${totalLinks}`);
  console.log(`   Broken links: ${brokenLinks}`);
  
  if (issues.length > 0) {
    console.log('\n❌ Broken links found:');
    for (const issue of issues) {
      console.log(`   ${issue.file}:${issue.line} - "${issue.text}" -> ${issue.url}`);
    }
    console.log('\n💡 Tip: Make sure all referenced files exist and paths are correct.');
    process.exit(1);
  } else {
    console.log('\n✅ All internal links are valid!');
  }
}

checkLinks();
