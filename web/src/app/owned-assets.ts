export class OwnedAssets {
    address: string
    assetIds: string[]

    constructor(init?: Partial<OwnedAssets>) {
        Object.assign(this, init)
    }
}
