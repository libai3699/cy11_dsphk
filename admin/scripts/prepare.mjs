import { spawnSync } from 'node:child_process';

if (process.env.CI) {
  process.exit(0);
}

const gitCheck = spawnSync('git', ['rev-parse', '--is-inside-work-tree'], {
  stdio: 'ignore',
});

if (gitCheck.status !== 0) {
  console.log('prepare: skip lefthook install because this directory is not inside a git worktree');
  process.exit(0);
}

const result = spawnSync('pnpm', ['exec', 'lefthook', 'install'], {
  shell: true,
  stdio: 'inherit',
});

process.exit(result.status ?? 1);
