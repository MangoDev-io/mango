import { Component, OnInit } from '@angular/core';
import { Router }       from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})

export class LoginComponent implements OnInit {

  public mnemonic: string;

  constructor(
    private router: Router,
  ) { }

  ngOnInit() {
  }

  loginClickHandler() {
    if (!this.mnemonic) {
      alert('Please enter the mnemonic');
    } else {
      console.log(this.mnemonic);
      this.router.navigate(['/manage']);
    }
  }

}
