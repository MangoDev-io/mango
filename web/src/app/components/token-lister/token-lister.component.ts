import { Component, OnInit } from '@angular/core';
import { StateService } from '../../state.service';

@Component({
  selector: 'app-token-lister',
  templateUrl: './token-lister.component.html',
  styleUrls: ['./token-lister.component.scss']
})
export class TokenListerComponent implements OnInit {

  constructor(public stateSvc: StateService) { }

  ngOnInit(): void {
  }

}
