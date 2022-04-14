import * as models from './models';

export interface go {
  "main": {
    "App": {
		Alert(arg1:string,arg2:string):Promise<void>
		GetCurrentFile():Promise<string>
		New(arg1:Array<number>):Promise<boolean>
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
