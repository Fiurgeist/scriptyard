process.env.TZ = 'UTC';
module.exports = {
  transform: {
    '\\.tsx?$': 'ts-jest',
  },
  setupFilesAfterEnv: ['<rootDir>/tests/setup.ts'],
  verbose: true,
  clearMocks: true,
  collectCoverageFrom: ['App.tsx'],
  coverageDirectory: 'coverage',
  moduleFileExtensions: ['js', 'jsx', 'ts', 'tsx'],
  testEnvironment: 'jsdom',
  coverageReporters: ['text', 'cobertura'],
  rootDir: '../',
};
