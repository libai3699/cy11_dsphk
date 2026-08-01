// 配置项
import envDev from './env/dev.js';
import envProd from './env/prod.js';

const compileEnv = process.env.APP_ENV || 'prod';

// 环境配置字典
const config = {
  // 开发环境配置
  dev: envDev,
  // 生产环境配置
  prod: envProd
};

// 最终应用的环境 - 编译环境
let CONFIG = config[compileEnv];

// 公共配置
const configCommon = {
  // 默认分页大小
  pageSize: 10
};

// 合并公共配置
CONFIG = Object.assign({}, configCommon, CONFIG);

export default CONFIG;
export { compileEnv, CONFIG };
export const APP_ENV = process.env.APP_ENV || 'dev';

