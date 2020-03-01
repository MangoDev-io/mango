import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { StateService } from 'src/app/state.service'

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.scss'],
})
export class LoginComponent implements OnInit {
    public mnemonic: string

    constructor(private router: Router, private stateService: StateService) {}

    ngOnInit() {}

    loginClickHandler() {
        if (!this.mnemonic) {
            alert('Please enter the mnemonic')
        } else {
            this.stateService.setMnemonic(this.mnemonic)
            this.router.navigate(['/manage'])
        }
    }
}
