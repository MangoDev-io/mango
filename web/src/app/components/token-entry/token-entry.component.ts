import { Component, OnInit, Input } from '@angular/core'
import { Token } from 'src/app/model/token'

@Component({
    selector: 'app-token-entry',
    templateUrl: './token-entry.component.html',
    styleUrls: ['./token-entry.component.scss'],
})
export class TokenEntryComponent implements OnInit {
    @Input()
    public token: Token

    public gradient: string

    constructor() {}

    ngOnInit() {
        this.generateGradient();
    }

    private generateGradient(): void {
        // prettier-ignore
        var hexValues = ["0","1","2","3","4","5","6","7","8","9","a","b","c","d","e"];

        function populate(a) {
            for (var i = 0; i < 6; i++) {
                var x = Math.round(Math.random() * 14)
                var y = hexValues[x]
                a += y
            }
            return a
        }

        var newColor1 = populate('#')
        var newColor2 = populate('#')
        var angle = Math.round(Math.random() * 360)

        // prettier-ignore
        this.gradient = "linear-gradient(" + angle + "deg, " + newColor1 + ", " + newColor2 + ")";
    }
}
