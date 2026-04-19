export default function AuthLayout({ children }: { children: React.ReactNode }) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-50 to-slate-200">
        {children}
      </main>
    );
  }
  