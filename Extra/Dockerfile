FROM ubuntu:24.04

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

RUN wget https://github.com/godotengine/godot/releases/download/4.4-stable/Godot_v4.4-stable_linux.x86_64.zip -O /tmp/godot.zip \
    && unzip /tmp/godot.zip -d /tmp \
    && chmod +x /tmp/Godot_v4.4-stable_linux.x86_64 \
    && mv /tmp/Godot_v4.4-stable_linux.x86_64 /usr/local/bin/godot \
    && rm /tmp/godot.zip 

WORKDIR /godot_project

COPY websocket_multiplayer/ /godot_project/

RUN echo '#!/bin/bash\n\
echo "Compilando proyecto WebSocket Multiplayer de Godot..."\n\
cd /godot_project\n\
# Crear una configuración de exportación básica si no existe\n\
if [ ! -f /godot_project/export_presets.cfg ]; then\n\
  echo "Configurando preset de exportación..."\n\
  echo "[preset.0]\\n\\nname=\\"Linux/X11\\"\\nplatform=\\"Linux/X11\\"\\nexport_path=\\"/output/websocket_multiplayer.x86_64\\"\\n" > /godot_project/export_presets.cfg\n\
fi\n\
godot --headless --verbose --export-debug "Linux/X11" /output/websocket_multiplayer.x86_64\n\
echo "Compilación completada. El binario se encuentra en /output/websocket_multiplayer.x86_64"\n\
' > /compile.sh \
    && chmod +x /compile.sh

# Crear directorio de salida
RUN mkdir -p /output && chmod 77 /output 

# Ejecutar el script de compilación
CMD ["/compile.sh"]