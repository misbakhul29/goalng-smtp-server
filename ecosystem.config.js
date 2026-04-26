module.exports = {
  apps: [
    {
      name: "golangsmtp",
      script: "./golangsmtp.exe",
      interpreter: "none",
      watch: false,
      autorestart: true,
      restart_delay: 3000,
      max_restarts: 10,
      env: {
        SERVER_PORT: "8000",
      },
    },
  ],
};
