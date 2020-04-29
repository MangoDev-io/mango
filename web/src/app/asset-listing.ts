export class AssetListing {
    assetId: string

    constructor(init?: Partial<AssetListing>) {
        Object.assign(this, init)
    }
}
