import { motion } from 'motion/react';
import { Github, Activity, Bell, Server, Terminal, Shield, Zap } from 'lucide-react';

export default function App() {
  const features = [
    {
      icon: <Activity className="w-6 h-6 text-neon-cyan" />,
      title: "Continuous Monitoring",
      description: "Constantly tracks your internet connection quality, latency, and packet loss in real-time."
    },
    {
      icon: <Bell className="w-6 h-6 text-neon-purple" />,
      title: "Telegram Alerts",
      description: "Get instant notifications directly to your Telegram when performance drops below acceptable levels."
    },
    {
      icon: <Zap className="w-6 h-6 text-neon-cyan" />,
      title: "Golang Powered",
      description: "Built with Go for extreme performance, low memory footprint, and robust concurrency."
    },
    {
      icon: <Server className="w-6 h-6 text-neon-purple" />,
      title: "Orange Pi 5 Optimized",
      description: "Specifically tuned and tested for ARM64 architecture and Orange Pi 5 hardware."
    },
    {
      icon: <Shield className="w-6 h-6 text-neon-cyan" />,
      title: "Systemd Auto-start",
      description: "Includes native systemd service configuration for reliable background execution and auto-restart."
    },
    {
      icon: <Terminal className="w-6 h-6 text-neon-purple" />,
      title: "Kubernetes Ready",
      description: "Full K3s deployment support with Docker images and secret management out of the box."
    }
  ];

  return (
    <div className="min-h-screen bg-dark-bg text-white font-sans selection:bg-neon-cyan selection:text-black">
      {/* Background Grid */}
      <div className="fixed inset-0 bg-grid-pattern opacity-20 pointer-events-none z-0"></div>
      
      {/* Ambient Glows */}
      <div className="fixed top-[-20%] left-[-10%] w-[50%] h-[50%] rounded-full bg-neon-cyan opacity-10 blur-[120px] pointer-events-none z-0"></div>
      <div className="fixed bottom-[-20%] right-[-10%] w-[50%] h-[50%] rounded-full bg-neon-purple opacity-10 blur-[120px] pointer-events-none z-0"></div>

      {/* Navigation */}
      <nav className="relative z-10 flex items-center justify-between px-6 py-6 max-w-7xl mx-auto">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded bg-dark-surface border border-dark-border flex items-center justify-center glow-cyan">
            <Activity className="w-5 h-5 text-neon-cyan" />
          </div>
          <span className="font-bold text-xl tracking-wider">TETRA</span>
        </div>
        <a 
          href="https://github.com/PiterPentester/tetra_bot" 
          target="_blank" 
          rel="noopener noreferrer"
          className="flex items-center gap-2 px-4 py-2 rounded-md bg-dark-surface border border-dark-border hover:border-neon-cyan transition-colors duration-300"
        >
          <Github className="w-5 h-5" />
          <span className="font-medium text-sm uppercase tracking-widest">GitHub</span>
        </a>
      </nav>

      <main className="relative z-10 max-w-7xl mx-auto px-6 pt-20 pb-32">
        {/* Hero Section */}
        <section className="flex flex-col items-center text-center mb-32">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-dark-surface border border-dark-border mb-8"
          >
            <span className="w-2 h-2 rounded-full bg-neon-cyan animate-pulse"></span>
            <span className="text-xs font-mono text-gray-400 uppercase tracking-widest">System Online</span>
          </motion.div>
          
          <motion.h1 
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.8, delay: 0.1 }}
            className="text-6xl md:text-8xl font-bold tracking-tighter mb-6"
          >
            TETRA <span className="text-transparent bg-clip-text bg-gradient-to-r from-neon-cyan to-neon-purple">BOT</span>
          </motion.h1>
          
          <motion.p 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.2 }}
            className="text-lg md:text-xl text-gray-400 max-w-2xl mb-10 font-light"
          >
            A robust Golang-based Telegram bot designed to continuously monitor your internet connection quality and alert you when performance drops.
          </motion.p>
          
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.3 }}
            className="flex flex-col sm:flex-row gap-4"
          >
            <a 
              href="https://github.com/PiterPentester/tetra_bot"
              target="_blank"
              rel="noopener noreferrer"
              className="px-8 py-4 rounded bg-white text-black font-bold uppercase tracking-widest hover:bg-neon-cyan transition-colors duration-300 glow-cyan flex items-center justify-center gap-2"
            >
              <Github className="w-5 h-5" />
              View Source
            </a>
            <a 
              href="#getting-started"
              className="px-8 py-4 rounded bg-dark-surface border border-dark-border text-white font-bold uppercase tracking-widest hover:border-neon-purple transition-colors duration-300 flex items-center justify-center gap-2"
            >
              <Terminal className="w-5 h-5" />
              Quick Start
            </a>
          </motion.div>
        </section>

        {/* Features Grid */}
        <section className="mb-32">
          <div className="flex items-center gap-4 mb-12">
            <h2 className="text-3xl font-bold tracking-tight uppercase">Core Features</h2>
            <div className="h-[1px] flex-1 bg-gradient-to-r from-dark-border to-transparent"></div>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((feature, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ duration: 0.5, delay: index * 0.1 }}
                className="p-6 rounded-xl bg-dark-surface border border-dark-border hover:border-neon-cyan/50 transition-colors duration-300 group"
              >
                <div className="w-12 h-12 rounded bg-dark-bg border border-dark-border flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                  {feature.icon}
                </div>
                <h3 className="text-xl font-semibold mb-3">{feature.title}</h3>
                <p className="text-gray-400 text-sm leading-relaxed">{feature.description}</p>
              </motion.div>
            ))}
          </div>
        </section>

        {/* Getting Started */}
        <section id="getting-started" className="mb-20">
          <div className="flex items-center gap-4 mb-12">
            <h2 className="text-3xl font-bold tracking-tight uppercase">Deployment</h2>
            <div className="h-[1px] flex-1 bg-gradient-to-r from-dark-border to-transparent"></div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <motion.div 
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              className="p-8 rounded-xl bg-dark-surface border border-dark-border"
            >
              <h3 className="text-xl font-bold mb-6 flex items-center gap-3">
                <Terminal className="w-5 h-5 text-neon-cyan" />
                Standard Installation
              </h3>
              <div className="space-y-4">
                <div className="bg-dark-bg p-4 rounded border border-dark-border font-mono text-sm overflow-x-auto">
                  <p className="text-gray-500 mb-1"># Clone the repository</p>
                  <p className="text-neon-cyan">git clone https://github.com/PiterPentester/tetra_bot.git</p>
                  <p className="text-neon-cyan">cd tetra_bot</p>
                  <br />
                  <p className="text-gray-500 mb-1"># Build the binary</p>
                  <p className="text-neon-cyan">go build -o tetra main.go</p>
                  <br />
                  <p className="text-gray-500 mb-1"># Run the bot</p>
                  <p className="text-neon-cyan">./tetra</p>
                </div>
              </div>
            </motion.div>

            <motion.div 
              initial={{ opacity: 0, x: 20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              className="p-8 rounded-xl bg-dark-surface border border-dark-border"
            >
              <h3 className="text-xl font-bold mb-6 flex items-center gap-3">
                <Server className="w-5 h-5 text-neon-purple" />
                Kubernetes (K3s)
              </h3>
              <div className="space-y-4">
                <div className="bg-dark-bg p-4 rounded border border-dark-border font-mono text-sm overflow-x-auto">
                  <p className="text-gray-500 mb-1"># Create namespace & secret</p>
                  <p className="text-neon-purple">kubectl create namespace tetra</p>
                  <p className="text-neon-purple">kubectl create secret generic tetra-secrets \</p>
                  <p className="text-neon-purple">  --from-literal=TELEGRAM_BOT_TOKEN=your_token</p>
                  <br />
                  <p className="text-gray-500 mb-1"># Apply deployment</p>
                  <p className="text-neon-purple">kubectl apply -f k8s/deployment.yaml</p>
                </div>
              </div>
            </motion.div>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="relative z-10 border-t border-dark-border bg-dark-surface py-12">
        <div className="max-w-7xl mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-6">
          <div className="flex items-center gap-2">
            <Activity className="w-5 h-5 text-neon-cyan" />
            <span className="font-bold tracking-widest text-sm">TETRA BOT</span>
          </div>
          <p className="text-gray-500 text-sm font-mono">
            Developed by <a href="https://github.com/PiterPentester" target="_blank" rel="noopener noreferrer" className="text-neon-cyan hover:underline">PiterPentester</a>
          </p>
          <div className="flex items-center gap-4">
            <a href="https://github.com/PiterPentester/tetra_bot" target="_blank" rel="noopener noreferrer" className="text-gray-400 hover:text-white transition-colors">
              <Github className="w-5 h-5" />
            </a>
          </div>
        </div>
      </footer>
    </div>
  );
}
