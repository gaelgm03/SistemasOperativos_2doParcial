FROM ubuntu:24.04

# Instalar dependencias necesarias
RUN apt-get update && apt-get install -y \
    wget \
    unzip \
    libxcursor1 \
    libxinerama1 \
    libxrandr2 \
    libxi6 \
    libgl1 \
    libsndio7.0 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Descargar y configurar Godot
RUN wget https://github.com/godotengine/godot/releases/download/4.4-stable/Godot_v4.4-stable_linux.x86_64.zip -O /tmp/godot.zip \
    && unzip /tmp/godot.zip -d /tmp \
    && chmod +x /tmp/Godot_v4.4-stable_linux.x86_64 \
    && mv /tmp/Godot_v4.4-stable_linux.x86_64 /usr/local/bin/godot \
    && rm /tmp/godot.zip 

# Crear un directorio de trabajo para el proyecto
WORKDIR /godot_project

# Ejecutar Godot en modo sin interfaz
CMD ["godot", "--headless"]