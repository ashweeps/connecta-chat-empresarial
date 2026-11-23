
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from argon2 import PasswordHasher
from argon2.exceptions import VerifyMismatchError, VerificationError

app = FastAPI(title="CONNECTA - Autenticacion")
hasher = PasswordHasher()

class PasswordRequest(BaseModel):
    password: str

class VerifyRequest(BaseModel):
    password: str
    password_hash: str

@app.get("/salud")
def salud():
    return {"estado": "activo"}

@app.post("/auth/hash")
def generar_hash(datos: PasswordRequest):
    if len(datos.password) < 12:
        raise HTTPException(
            status_code=400,
            detail="La contraseña debe tener al menos 12 caracteres"
        )

    return {"password_hash": hasher.hash(datos.password)}

@app.post("/auth/verificar")
def verificar_password(datos: VerifyRequest):
    try:
        valido = hasher.verify(
            datos.password_hash,
            datos.password
        )
    except (VerifyMismatchError, VerificationError, ValueError):
        valido = False

    if not valido:
        raise HTTPException(
            status_code=401,
            detail="Credenciales incorrectas"
        )

    return {"valido": True}
