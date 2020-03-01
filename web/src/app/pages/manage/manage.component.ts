import { Component, OnInit } from '@angular/core'

@Component({
    selector: 'app-manage',
    templateUrl: './manage.component.html',
    styleUrls: ['./manage.component.scss'],
})
export class ManageComponent implements OnInit {
    public showTokenCreate = true

    constructor() {}

    ngOnInit() {}
}
