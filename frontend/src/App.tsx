import React, { useState } from 'react';
import { Login } from './components/Login/Login';
import { PostList } from './components/PostList/PostList';
import { CreatePost } from './components/CreatePost/CreatePost';
import { PostDetail } from './components/PostDetail/PostDetail';
import { User } from './types';
import './App.css';

type View = 'list' | 'detail';

function App() {
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [refreshPosts, setRefreshPosts] = useState(false);
  const [currentView, setCurrentView] = useState<View>('list');
  const [selectedPostId, setSelectedPostId] = useState<number | null>(null);

  const handleLoginSuccess = (user: User) => {
    setCurrentUser(user);
  };

  const handleLogout = () => {
    setCurrentUser(null);
    setCurrentView('list');
    setSelectedPostId(null);
  };

  const handlePostCreated = () => {
    setRefreshPosts(!refreshPosts);
  };

  const handleViewPost = (postId: number) => {
    setSelectedPostId(postId);
    setCurrentView('detail');
  };

  const handleBackToList = () => {
    setCurrentView('list');
    setSelectedPostId(null);
    setRefreshPosts(!refreshPosts);
  };

  // Si no está logueado, mostrar login
  if (!currentUser) {
    return <Login onLoginSuccess={handleLoginSuccess} />;
  }

  // Si está logueado, mostrar la app
  return (
    <div className="App">
      <header className="app-header">
        <h1>🚀 Mini Red Social</h1>
        <div className="user-info">
          <span>Hola, @{currentUser.username}</span>
          <button onClick={handleLogout} className="logout-btn">
            Cerrar Sesión
          </button>
        </div>
      </header>

      <main>
        {currentView === 'list' ? (
          <>
            <CreatePost userId={currentUser.id} onPostCreated={handlePostCreated} />
            <PostList 
              currentUserId={currentUser.id} 
              onRefresh={refreshPosts}
              onViewPost={handleViewPost}
            />
          </>
        ) : (
          selectedPostId && (
            <PostDetail 
              postId={selectedPostId}
              userId={currentUser.id}
              onBack={handleBackToList}
            />
          )
        )}
      </main>
    </div>
  );
}

export default App;