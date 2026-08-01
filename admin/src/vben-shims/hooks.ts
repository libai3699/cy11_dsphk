// @vben/hooks 替代实现

export function useRequest() {
  return {
    loading: false,
    run: async (fn: Function) => fn(),
  };
}
