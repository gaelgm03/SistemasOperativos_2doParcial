import React, { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [tasks, setTasks] = useState([]);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [editingTask, setEditingTask] = useState(null);

  const API_URL = '/api';

  useEffect(() => {
    fetchTasks();
  }, []);

  const fetchTasks = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_URL}/tasks`);
      if (!response.ok) {
        throw new Error('Error al cargar las tareas');
      }
      const data = await response.json();
      setTasks(data);
      setError(null);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!title.trim()) {
      alert('El título es obligatorio');
      return;
    }

    try {
      if (editingTask) {
        // Actualizar tarea existente
        await fetch(`${API_URL}/tasks/${editingTask.id}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            title,
            description,
            completed: editingTask.completed
          })
        });
        setEditingTask(null);
      } else {
        // Crear nueva tarea
        await fetch(`${API_URL}/tasks`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            title,
            description,
            completed: false
          })
        });
      }
      
      // Resetear formulario y recargar tareas
      setTitle('');
      setDescription('');
      fetchTasks();
    } catch (err) {
      setError(err.message);
    }
  };

  const toggleComplete = async (task) => {
    try {
      await fetch(`${API_URL}/tasks/${task.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          ...task,
          completed: !task.completed
        })
      });
      fetchTasks();
    } catch (err) {
      setError(err.message);
    }
  };

  const deleteTask = async (id) => {
    if (!window.confirm('¿Estás seguro de que quieres eliminar esta tarea?')) {
      return;
    }
    
    try {
      await fetch(`${API_URL}/tasks/${id}`, {
        method: 'DELETE'
      });
      fetchTasks();
    } catch (err) {
      setError(err.message);
    }
  };

  const startEdit = (task) => {
    setEditingTask(task);
    setTitle(task.title);
    setDescription(task.description || '');
  };

  const cancelEdit = () => {
    setEditingTask(null);
    setTitle('');
    setDescription('');
  };

  return (
    <div className="app">
      <header>
        <h1>Administrador de Tareas</h1>
      </header>
      
      <main>
        <section className="form-section">
          <h2>{editingTask ? 'Editar Tarea' : 'Crear Nueva Tarea'}</h2>
          
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="title">Título:</label>
              <input
                type="text"
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Título de la tarea"
                required
              />
            </div>
            
            <div className="form-group">
              <label htmlFor="description">Descripción:</label>
              <textarea
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Descripción (opcional)"
                rows="3"
              />
            </div>
            
            <div className="form-buttons">
              <button type="submit" className="btn primary">
                {editingTask ? 'Actualizar' : 'Crear'}
              </button>
              
              {editingTask && (
                <button 
                  type="button" 
                  className="btn secondary"
                  onClick={cancelEdit}
                >
                  Cancelar
                </button>
              )}
            </div>
          </form>
        </section>
        
        <section className="tasks-section">
          <h2>Mis Tareas</h2>
          
          {loading && <p className="loading">Cargando tareas...</p>}
          
          {error && (
            <div className="error-message">
              Error: {error}
              <button onClick={fetchTasks} className="btn">Reintentar</button>
            </div>
          )}
          
          {!loading && !error && tasks.length === 0 && (
            <p className="no-tasks">No hay tareas disponibles.</p>
          )}
          
          <ul className="task-list">
            {tasks.map((task) => (
              <li 
                key={task.id} 
                className={`task-item ${task.completed ? 'completed' : ''}`}
              >
                <div className="task-content">
                  <h3>{task.title}</h3>
                  {task.description && <p>{task.description}</p>}
                  <span className="created-date">
                    Creada: {new Date(task.created_at).toLocaleDateString()}
                  </span>
                </div>
                
                <div className="task-actions">
                  <button 
                    onClick={() => toggleComplete(task)}
                    className="btn-icon"
                    title={task.completed ? "Marcar como pendiente" : "Marcar como completada"}
                  >
                    {task.completed ? '↩️' : '✅'}
                  </button>
                  
                  <button 
                    onClick={() => startEdit(task)}
                    className="btn-icon"
                    title="Editar"
                  >
                    ✏️
                  </button>
                  
                  <button 
                    onClick={() => deleteTask(task.id)}
                    className="btn-icon"
                    title="Eliminar"
                  >
                    🗑️
                  </button>
                </div>
              </li>
            ))}
          </ul>
        </section>
      </main>
      
      <footer>
        <p>API con Go y Frontend con React - © 2025</p>
      </footer>
    </div>
  );
}

export default App;