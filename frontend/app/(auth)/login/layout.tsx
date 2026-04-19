export const metadata = {
    title: 'Login - Flashcards',
    description: 'Login to Flashcards and start mastering languages with ease.',
  };
  
  export default function LoginLayout({ children }: { children: React.ReactNode }) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-indigo-100 to-white p-4">
        {children}
      </div>
    );
  }
  