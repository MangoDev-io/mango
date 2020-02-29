import { Injectable } from '@angular/core';
import { Token } from './model/token';

@Injectable({
  providedIn: 'root'
})
export class StateService {
  public tokenList: Token[] = [
    { id: 1, name: 'A'},
    { id: 2, name: 'B'},
    { id: 3, name: 'C'}
  ]

  constructor() { }
}
