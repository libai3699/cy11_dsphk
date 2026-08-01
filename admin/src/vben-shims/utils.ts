// @vben/utils 替代实现

export function unmountGlobalLoading() {
  const loadingEl = document.getElementById('loading-container');
  if (loadingEl) {
    loadingEl.style.display = 'none';
  }
}

export function mergeRouteModules(modules: Record<string, any>) {
  const routes: any[] = [];
  Object.keys(modules).forEach((key) => {
    const mod = modules[key].default || modules[key];
    const modList = Array.isArray(mod) ? [...mod] : [mod];
    routes.push(...modList);
  });
  return routes;
}

export function traverseTreeValues(tree: any[], callback: (item: any) => any) {
  const result: any[] = [];
  function traverse(nodes: any[]) {
    nodes.forEach((node) => {
      result.push(callback(node));
      if (node.children) {
        traverse(node.children);
      }
    });
  }
  traverse(tree);
  return result;
}

export function openWindow(url: string, options?: any) {
  window.open(url, options?.target || '_blank');
}

export function startProgress() {
  // 简化实现
}

export function stopProgress() {
  // 简化实现
}
