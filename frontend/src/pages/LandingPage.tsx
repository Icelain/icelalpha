import { createSignal } from "solid-js";
import { FiGithub } from 'solid-icons/fi';

function LandingPage() {
  const [menuOpen, setMenuOpen] = createSignal(false);

  return (
    <div class="min-h-screen bg-black text-green-400 font-mono flex flex-col">
      {/* Minimal Header */}
      <header class="border-b border-green-800 py-4">
        <div class="container mx-auto px-4 flex items-center justify-between">
          <h1 class="text-xl md:text-2xl font-bold tracking-tight">
            <span class="text-green-500">&gt;_</span> IceAlpha
          </h1>

          {/* Desktop Nav */}
          <nav class="hidden md:flex items-center space-x-8">
            <a href="#about" class="hover:text-green-300 transition-colors">/about</a>
            <a href="#examples" class="hover:text-green-300 transition-colors">/examples</a>
            <button 
              class="bg-green-900 hover:bg-green-800 text-green-300 rounded px-4 py-2 flex items-center space-x-2 transition-colors"
              onclick={() => window.location.href = '/auth/github'}
            >
              <FiGithub />
              <span>github_login</span>
            </button>
          </nav>

          {/* Mobile menu button */}
          <button 
            class="md:hidden text-green-400 p-2" 
            onclick={() => setMenuOpen(!menuOpen())}
            aria-label="Toggle menu"
          >
            <span class="text-lg">{menuOpen() ? 'x' : '≡'}</span>
          </button>
        </div>

        {/* Mobile Navigation */}
        <nav 
          class={`md:hidden bg-black ${menuOpen() ? 'block' : 'hidden'}`}
        >
          <div class="container mx-auto px-4 py-3 flex flex-col space-y-3">
            <a href="#about" class="block py-2 hover:text-green-300 transition-colors">/about</a>
            <a href="#examples" class="block py-2 hover:text-green-300 transition-colors">/examples</a>
            <button 
              class="bg-green-900 hover:bg-green-800 text-green-300 rounded px-4 py-2 flex items-center space-x-2 transition-colors w-full"
              onclick={() => window.location.href = '/auth/github'}
            >
              <FiGithub />
              <span>github_login</span>
            </button>
          </div>
        </nav>
      </header>

      {/* Minimal Terminal-like Hero */}
      <section class="py-12 px-4 flex-grow flex items-center">
        <div class="container mx-auto max-w-3xl">
          <div class="bg-gray-900 rounded border border-green-800 p-4 md:p-6 shadow-lg">
            <div class="flex items-center mb-2">
              <div class="h-3 w-3 rounded-full bg-red-500 mr-2"></div>
              <div class="h-3 w-3 rounded-full bg-yellow-500 mr-2"></div>
              <div class="h-3 w-3 rounded-full bg-green-500"></div>
              <span class="ml-4 text-xs text-gray-400">icealpha@math:~</span>
            </div>
            <div class="space-y-4">
              <p class="text-sm md:text-base">
                <span class="text-green-500">$</span> <span class="text-blue-400">whatis</span> icealpha
              </p>
              <p class="text-sm md:text-base pl-4">
                IceAlpha: AI-powered mathematical problem solver
              </p>
              <p class="text-sm md:text-base">
                <span class="text-green-500">$</span> <span class="text-blue-400">cat</span> features.txt
              </p>
              <div class="pl-4">
                <p class="text-sm md:text-base">- Advanced symbolic computation</p>
                <p class="text-sm md:text-base">- Step-by-step solution explanations</p>
                <p class="text-sm md:text-base">- LaTeX integration for formatted math</p>
                <p class="text-sm md:text-base">- Natural language query processing</p>
              </div>
              <p class="text-sm md:text-base">
                <span class="text-green-500">$</span> <span class="text-blue-400">./</span> connect_github.sh
              </p>
              <p>
                <button 
                  class="bg-green-900 hover:bg-green-800 text-green-300 rounded px-4 py-2 flex items-center space-x-2 transition-colors"
                  onclick={() => window.location.href = '/auth/github'}
                >
                  <FiGithub />
                  <span>Authenticate with GitHub</span>
                </button>
              </p>
              <p class="text-sm md:text-base flex items-center">
                <span class="text-green-500 mr-2">$</span> <span class="animate-pulse">█</span>
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Examples Section */}
      <section id="examples" class="py-8 px-4 bg-gray-900">
        <div class="container mx-auto max-w-3xl">
          <h2 class="text-xl border-b border-green-800 pb-2 mb-6">> examples/</h2>
          <div class="bg-black p-4 rounded border border-green-800 mb-6">
            <p class="text-gray-400 mb-2">// Solving a differential equation</p>
            <p class="mb-2"><span class="text-blue-400">input</span>: solve dy/dx = 2x + y</p>
            <p class="text-green-400">y = Ce^x - 2x - 2</p>
          </div>
          <div class="bg-black p-4 rounded border border-green-800">
            <p class="text-gray-400 mb-2">// Calculating a complex integral</p>
            <p class="mb-2"><span class="text-blue-400">input</span>: integrate x^2 * sin(x) from 0 to pi</p>
            <p class="text-green-400">2π - 4</p>
          </div>
        </div>
      </section>

      {/* About Section */}
      <section id="about" class="py-8 px-4">
        <div class="container mx-auto max-w-3xl">
          <h2 class="text-xl border-b border-green-800 pb-2 mb-6">> about/</h2>
          <pre class="bg-black p-4 rounded border border-green-800 text-sm md:text-base whitespace-pre-wrap">
{`IceAlpha is an AI-powered tool designed to solve 
complex mathematical problems through advanced 
symbolic computation and machine learning.

Built for students, researchers, and professionals 
who need fast, accurate solutions with detailed 
explanation steps.

Version: 1.2.4
License: MIT
Author: IceAlpha Team`}
          </pre>
        </div>
      </section>

      {/* Footer */}
      <footer class="border-t border-green-800 py-4 px-4 mt-auto">
        <div class="container mx-auto max-w-3xl flex flex-col md:flex-row justify-between items-center">
          <p class="text-xs text-gray-500 mb-4 md:mb-0">© 2025 IceAlpha | [v1.2.4]</p>
          <div class="flex space-x-6">
            <a href="#" class="text-xs text-gray-500 hover:text-green-400 transition-colors">
              GitHub
            </a>
            <a href="#" class="text-xs text-gray-500 hover:text-green-400 transition-colors">
              Docs
            </a>
            <a href="#" class="text-xs text-gray-500 hover:text-green-400 transition-colors">
              API
            </a>
          </div>
        </div>
      </footer>
    </div>
  );
}

export default LandingPage;