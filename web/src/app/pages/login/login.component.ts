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

    showNotificationModal = false
    responseError = 'Invalid Mnemonic'

    constructor(private router: Router, private stateService: StateService) {}

    ngOnInit() {}

    loginClickHandler() {
        if (!this.mnemonic) {
            alert('Please enter the mnemonic')
        } else {
            this.stateService.encodeMnemonic(this.mnemonic).subscribe(
                (x: any) => {
                    this.stateService.setAuthToken(x.token)
                    this.router.navigate(['/manage'])
                },
                err => {
                    this.showNotificationModal = true
                    this.showNotificationModal = false
                    this.showNotificationModal = true
                    console.error(err)
                }
            )
        }
    }
}
