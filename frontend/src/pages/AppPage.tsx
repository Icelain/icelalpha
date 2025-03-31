
import { createSignal, createEffect, Show } from "solid-js";
import { FiGithub, FiLogOut, FiImage, FiCode, FiSend, FiX } from 'solid-icons/fi';

function AppPage() {
  const [inputMode, setInputMode] = createSignal("text"); // "text", "latex", "image"
  const [queryText, setQueryText] = createSignal("");
  const [latexCode, setLatexCode] = createSignal("");
  const [imageUploaded, setImageUploaded] = createSignal(false);
  const [isSubmitting, setIsSubmitting] = createSignal(false);
  const [result, setResult] = createSignal(null);
  const [history, setHistory] = createSignal([]);

  const handleSubmit = (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    
    // Simulate API call
    setTimeout(() => {
      const query = inputMode() === "text" ? queryText() : 
                    inputMode() === "latex" ? latexCode() : "Image query";
      
      // Example response - would come from your actual API
      const sampleResponse = {
        query: query,
        solution: "y = 2x² + C",
        steps: [
          "First, identify this as a separable differential equation",
          "Rearrange to get dy/dx = f(x) · g(y)",
          "Integrate both sides",
          "Apply the initial conditions",
          "Final solution: y = 2x² + C"
        ]
      };
      
      setResult(sampleResponse);
      setHistory([sampleResponse, ...history()]);
      setIsSubmitting(false);
      
      // Reset inputs
      setQueryText("");
      setLatexCode("");
      setImageUploaded(false);
    }, 1000);
  };

  const handleFileUpload = (e) => {
    const file = e.target.files[0];
    if (file) {
      // Normally you'd process the image here
      setImageUploaded(true);
    }
  };

  const clearResults = () => {
    setResult(null);
  };

  return (
    <div class="min-h-screen bg-black text-green-400 font-mono flex flex-col">
      {/* Header */}
      <header class="border-b border-green-800 py-4">
        <div class="container mx-auto px-4 flex items-center justify-between">
          <h1 class="text-xl md:text-2xl font-bold tracking-tight">
            <span class="text-green-500">&gt;_</span> IceAlpha
          </h1>

          <div class="flex items-center">
            <div class="mr-4 text-xs hidden md:block">
              user@icealpha
            </div>
            <button
              class="bg-green-900 hover:bg-green-800 text-green-300 rounded px-3 py-1 flex items-center space-x-2 transition-colors"
              onclick={() => window.location.href = '/logout'}
            >
              <FiLogOut size={14} />
              <span class="text-sm hidden md:inline">logout</span>
            </button>
          </div>
        </div>
      </header>

      <div class="flex flex-col md:flex-row flex-1">
        {/* Main Content */}
        <main class="flex-1 p-4">
          <div class="mb-6 bg-gray-900 rounded border border-green-800 p-4 shadow-lg">
            <div class="flex items-center mb-2">
              <div class="h-3 w-3 rounded-full bg-red-500 mr-2"></div>
              <div class="h-3 w-3 rounded-full bg-yellow-500 mr-2"></div>
              <div class="h-3 w-3 rounded-full bg-green-500"></div>
              <span class="ml-4 text-xs text-gray-400">math-solver:~</span>
            </div>
            
            <form onSubmit={handleSubmit}>
              <div class="space-y-4">
                {/* Input Type Selector */}
                <div class="flex space-x-2">
                  <button 
                    type="button"
                    class={`px-3 py-1 rounded text-sm ${inputMode() === "text" ? "bg-green-800 text-green-200" : "bg-gray-800 text-gray-400"}`}
                    onClick={() => setInputMode("text")}
                  >
                    Text
                  </button>
                  <button 
                    type="button"
                    class={`px-3 py-1 rounded text-sm flex items-center ${inputMode() === "latex" ? "bg-green-800 text-green-200" : "bg-gray-800 text-gray-400"}`}
                    onClick={() => setInputMode("latex")}
                  >
                    <FiCode class="mr-1" />
                    LaTeX
                  </button>
                  <button 
                    type="button"
                    class={`px-3 py-1 rounded text-sm flex items-center ${inputMode() === "image" ? "bg-green-800 text-green-200" : "bg-gray-800 text-gray-400"}`}
                    onClick={() => setInputMode("image")}
                  >
                    <FiImage class="mr-1" />
                    Image
                  </button>
                </div>
                
                {/* Input Fields - shown conditionally based on input mode */}
                <Show when={inputMode() === "text"}>
                  <div>
                    <p class="text-sm mb-2">
                      <span class="text-green-500">$</span> <span class="text-blue-400">solve</span>
                    </p>
                    <textarea
                      class="w-full bg-black text-green-300 border border-green-800 rounded p-2 focus:outline-none focus:border-green-600"
                      placeholder="Enter your mathematical problem here..."
                      rows="3"
                      value={queryText()}
                      onInput={(e) => setQueryText(e.target.value)}
                    ></textarea>
                  </div>
                </Show>
                
                <Show when={inputMode() === "latex"}>
                  <div>
                    <p class="text-sm mb-2">
                      <span class="text-green-500">$</span> <span class="text-blue-400">solve-latex</span>
                    </p>
                    <textarea
                      class="w-full bg-black text-green-300 border border-green-800 rounded p-2 focus:outline-none focus:border-green-600"
                      placeholder="\int_0^{\pi} x^2 \sin(x) dx"
                      rows="3"
                      value={latexCode()}
                      onInput={(e) => setLatexCode(e.target.value)}
                    ></textarea>
                  </div>
                </Show>
                
                <Show when={inputMode() === "image"}>
                  <div>
                    <p class="text-sm mb-2">
                      <span class="text-green-500">$</span> <span class="text-blue-400">solve-image</span>
                    </p>
                    <div class="border border-dashed border-green-800 rounded p-6 text-center">
                      {!imageUploaded() ? (
                        <>
                          <p class="mb-2 text-gray-400">Upload image of your math problem</p>
                          <input 
                            type="file" 
                            accept="image/*" 
                            class="hidden" 
                            id="file-upload" 
                            onChange={handleFileUpload}
                          />
                          <label 
                            for="file-upload" 
                            class="cursor-pointer bg-gray-800 hover:bg-gray-700 text-green-300 px-3 py-2 rounded inline-flex items-center"
                          >
                            <FiImage class="mr-2" />
                            Select Image
                          </label>
                        </>
                      ) : (
                        <div class="flex items-center justify-center">
                          <span class="text-green-400 mr-2">Image uploaded</span>
                          <button 
                            type="button" 
                            class="text-red-400 hover:text-red-300"
                            onClick={() => setImageUploaded(false)}
                          >
                            <FiX />
                          </button>
                        </div>
                      )}
                    </div>
                  </div>
                </Show>
                
                {/* Submit Button */}
                <div class="flex justify-end">
                  <button
                    type="submit"
                    disabled={isSubmitting() || 
                      (inputMode() === "text" && !queryText()) || 
                      (inputMode() === "latex" && !latexCode()) ||
                      (inputMode() === "image" && !imageUploaded())}
                    class="px-4 py-2 bg-green-900 hover:bg-green-800 text-green-300 rounded flex items-center space-x-2 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {isSubmitting() ? (
                      <span>Processing...</span>
                    ) : (
                      <>
                        <FiSend />
                        <span>Execute</span>
                      </>
                    )}
                  </button>
                </div>
              </div>
            </form>
          </div>
          
          {/* Results Section */}
          <Show when={result()}>
            <div class="bg-gray-900 rounded border border-green-800 p-4 shadow-lg mb-6">
              <div class="flex justify-between items-start mb-3">
                <h2 class="text-lg text-green-400">Results</h2>
                <button 
                  onClick={clearResults}
                  class="text-gray-400 hover:text-gray-300"
                >
                  <FiX />
                </button>
              </div>
              <div class="mb-3 pb-3 border-b border-gray-800">
                <span class="text-gray-400">Query:</span>
                <p class="text-green-300 mt-1">{result()?.query}</p>
              </div>
              <div class="mb-3 pb-3 border-b border-gray-800">
                <span class="text-gray-400">Solution:</span>
                <p class="text-green-300 mt-1 text-xl">{result()?.solution}</p>
              </div>
              <div>
                <span class="text-gray-400">Steps:</span>
                <ul class="mt-1 space-y-1">
                  {result()?.steps.map((step, index) => (
                    <li class="text-green-300">
                      <span class="text-gray-500">{index + 1}.</span> {step}
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </Show>
        </main>
        
        {/* History Sidebar */}
        <aside class="w-full md:w-64 p-4 border-t md:border-t-0 md:border-l border-green-800 overflow-y-auto">
          <h2 class="text-lg mb-4 pb-1 border-b border-green-800">History</h2>
          <div class="space-y-3">
            {history().length > 0 ? (
              history().map((item, index) => (
                <div class="p-2 border border-green-900 bg-gray-900 rounded text-xs hover:bg-gray-800">
                  <p class="text-gray-400 truncate">{item.query}</p>
                  <p class="text-green-400 truncate font-bold">{item.solution}</p>
                </div>
              ))
            ) : (
              <p class="text-gray-500 text-sm">No previous queries</p>
            )}
          </div>
        </aside>
      </div>
      
      {/* Footer */}
      <footer class="border-t border-green-800 py-3 px-4">
        <div class="container mx-auto text-center text-xs text-gray-500">
          IceAlpha v1.2.4 | AI-powered mathematical problem solver | © 2025
        </div>
      </footer>
    </div>
  );
}

export default AppPage;