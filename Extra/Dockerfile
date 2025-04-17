FROM ubuntu:24.04

# Evitar prompts durante la instalación
ENV DEBIAN_FRONTEND=noninteractive

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
    xvfb \
    libpulse0 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Descargar y configurar Godot
RUN wget https://github.com/godotengine/godot/releases/download/4.4-stable/Godot_v4.4-stable_linux.x86_64.zip -O /tmp/godot.zip \
    && unzip /tmp/godot.zip -d /tmp \
    && chmod +x /tmp/Godot_v4.4-stable_linux.x86_64 \
    && mv /tmp/Godot_v4.4-stable_linux.x86_64 /usr/local/bin/godot \
    && rm /tmp/godot.zip

# Descargar la plantilla de exportación para Linux
RUN mkdir -p /root/.local/share/godot/export_templates/4.4.stable && \
    wget https://github.com/godotengine/godot/releases/download/4.4-stable/Godot_v4.4-stable_export_templates.tpz -O /tmp/templates.tpz && \
    unzip /tmp/templates.tpz -d /tmp && \
    mv /tmp/templates/* /root/.local/share/godot/export_templates/4.4.stable/ && \
    rm -rf /tmp/templates.tpz /tmp/templates

# Crear directorios para el proyecto
WORKDIR /app

# Script de compilación
RUN echo '#!/bin/bash\n\
echo "Iniciando compilación del proyecto Godot..."\n\
\n\
# Verificar la existencia del directorio del proyecto\n\
if [ ! -d "/app/websocket_multiplayer" ]; then\n\
  echo "ERROR: No se encontró el directorio websocket_multiplayer."\n\
  echo "Contenido del directorio actual:"\n\
  ls -la /app\n\
  exit 1\n\
fi\n\
\n\
# Entrar al directorio del proyecto\n\
cd /app/websocket_multiplayer\n\
\n\
# Verificar si existe el archivo project.godot\n\
if [ ! -f "project.godot" ]; then\n\
  echo "ERROR: No se encontró el archivo project.godot en el directorio del proyecto."\n\
  echo "Contenido del directorio:"\n\
  ls -la\n\
  exit 1\n\
fi\n\
\n\
# Crear directorio de salida\n\
mkdir -p /output\n\
chmod 777 /output\n\
\n\
# Crear configuración de exportación si no existe\n\
if [ ! -f "export_presets.cfg" ]; then\n\
  echo "Creando configuración de exportación..."\n\
  cat > export_presets.cfg << EOL\n\
[preset.0]\n\
\n\
name="Linux/X11"\n\
platform="Linux/X11"\n\
runnable=true\n\
custom_features=""\n\
export_filter="all_resources"\n\
include_filter=""\n\
exclude_filter=""\n\
export_path="/output/websocket_game.x86_64"\n\
encryption_include_filters=""\n\
encryption_exclude_filters=""\n\
encrypt_pck=false\n\
encrypt_directory=false\n\
script_export_mode=1\n\
\n\
[preset.0.options]\n\
\n\
custom_template/debug=""\n\
custom_template/release=""\n\
binary_format/64_bits=true\n\
binary_format/embed_pck=false\n\
texture_format/bptc=false\n\
texture_format/s3tc=true\n\
texture_format/etc=false\n\
texture_format/etc2=false\n\
texture_format/no_bptc_fallbacks=true\n\
EOL\n\
fi\n\
\n\
# Mostrar información de diagnóstico\n\
echo "Versión de Godot:"\n\
godot --version\n\
\n\
echo "Contenido del directorio del proyecto:"\n\
ls -la\n\
\n\
echo "Verificando plantillas de exportación:"\n\
ls -la /root/.local/share/godot/export_templates/4.4.stable/\n\
\n\
# Realizar la exportación\n\
echo "Iniciando exportación..."\n\
xvfb-run --auto-servernum godot --headless --path /app/websocket_multiplayer --export-debug "Linux/X11" /output/websocket_game.x86_64 --verbose\n\
\n\
# Verificar el resultado\n\
if [ -f "/output/websocket_game.x86_64" ]; then\n\
  echo "Compilación EXITOSA. El binario se encuentra en /output/websocket_game.x86_64"\n\
  chmod 755 /output/websocket_game.x86_64\n\
  ls -la /output\n\
else\n\
  echo "Primer intento de exportación falló. Intentando método alternativo..."\n\
  # Intento alternativo\n\
  xvfb-run --auto-servernum godot --headless --path /app/websocket_multiplayer --export "Linux/X11" /output/websocket_game.x86_64 --verbose\n\
  \n\
  if [ -f "/output/websocket_game.x86_64" ]; then\n\
    echo "Compilación EXITOSA con método alternativo. El binario se encuentra en /output/websocket_game.x86_64"\n\
    chmod 755 /output/websocket_game.x86_64\n\
    ls -la /output\n\
  else\n\
    echo "ERROR: Todos los intentos de exportación fallaron."\n\
    echo "Contenido del directorio de salida:"\n\
    ls -la /output\n\
    echo "Información de diagnóstico adicional:"\n\
    godot --headless --path /app/websocket_multiplayer --verbose --doctor\n\
  fi\n\
fi\n\
' > /compile.sh \
    && chmod +x /compile.sh

CMD ["/compile.sh"]