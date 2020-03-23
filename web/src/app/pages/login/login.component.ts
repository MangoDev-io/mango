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
    notificationType = 1
    responseError: string

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
                    this.responseError = 'Invalid Mnemonic'
                    this.showNotificationModal = true
                    console.error(err)
                }
            )
        }
    }

    modalClosed() {
        this.showNotificationModal = false
    }
}
