import * as models from './models';

export interface go {
  "main": {
    "App": {
		GetCurrentFile():Promise<string>
		New():Promise<boolean>
		Open():Promise<boolean>
		Save():Promise<boolean>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
