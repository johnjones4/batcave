export const blobToBase64 = (blob: Blob): Promise<string> => {
  return new Promise((resolve, _) => {
    const reader = new FileReader()
    reader.onloadend = () => resolve(cleanBase64URL(reader.result as string))
    reader.readAsDataURL(blob)
  })
}

const cleanBase64URL = (u: string): string => {
  const start = u.indexOf(',')
  return u.substring(start + 1)
}