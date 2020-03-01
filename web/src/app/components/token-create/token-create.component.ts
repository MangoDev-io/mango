import { Component, OnInit } from '@angular/core'

@Component({
    selector: 'app-token-create',
    templateUrl: './token-create.component.html',
    styleUrls: ['./token-create.component.scss'],
})
export class TokenCreateComponent implements OnInit {
    public currentAddr =
        'BVMEUTF37WNEQ6GYCZISRFHGLEMOKT5OCPPTTJXVED6JBSXKF6YJJRZRI4'

    constructor() {}

    ngOnInit(): void {}
}
