export class Loader<T> {
  path: string

  constructor(path: string) {
    this.path = path
  }

  async load(): Promise<T> {
    const response = await fetch(this.path)
    const info = await response.json()
    return (info as T)
  }
}