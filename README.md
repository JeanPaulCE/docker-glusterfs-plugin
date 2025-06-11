# GlusterFS Volume Plugin

Plugin de volumen Docker para acceder a volúmenes GlusterFS. No requiere la instalación del cliente GlusterFS en el host.

## Requisitos

- Docker 18.03-1 o superior
- Plugin gestionado (managed plugin)

## Instalación

```bash
# Instalar el plugin
docker plugin install --alias glusterfs \
  mochoa/glusterfs-volume-plugin \
  --grant-all-permissions --disable

# Configurar servidores GlusterFS
docker plugin set glusterfs SERVERS=store1,store2

# Habilitar el plugin
docker plugin enable glusterfs
```

## Uso

### 1. Modo Simple (Recomendado)

```yaml
volumes:
  myvolume:
    driver: glusterfs
    name: "volume/subdir"
```

### 2. Modo con Servidores Específicos

```yaml
volumes:
  myvolume:
    driver: glusterfs
    driver_opts:
      servers: store1,store2
    name: "volume/subdir"
```

### 3. Modo con Opciones Personalizadas

```yaml
volumes:
  myvolume:
    driver: glusterfs
    driver_opts:
      glusteropts: "--volfile-server=SERVER --volfile-id=abc --subdir-mount=sub"
    name: "whatever"
```

## Ejemplo de Uso

```bash
# Crear un volumen
docker volume create -d glusterfs --opt servers=store1 myvol

# Usar el volumen en un contenedor
docker run -it -v myvol:/mnt alpine
```

## Configuración SSL

Para habilitar SSL en el canal de gestión:

```bash
docker plugin set glusterfs SECURE_MANAGEMENT=yes
```

## Notas Importantes

1. Los servidores GlusterFS deben estar definidos en `/etc/hosts` del runtime de Docker
2. Solo se soporta un cluster GlusterFS por instancia
3. Use `--alias` para definir instancias separadas si necesita diferentes configuraciones
