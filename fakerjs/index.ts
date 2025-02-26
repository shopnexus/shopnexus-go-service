import { faker } from "@faker-js/faker"
import {
	PrismaClient,
	Role,
	Gender,
	Brand,
	Account,
	ProductModel,
	Product,
} from "@prisma/client"
const prisma = new PrismaClient()

// Helper function to generate random enum values
const randomEnum = <T>(enumObject: { [key: string]: T }): T => {
	const values = Object.values(enumObject)
	return values[Math.floor(Math.random() * values.length)]
}

// Generate Account with UserAccount or AdminAccount
async function createAccounts(count: number) {
	const accounts: Account[] = []
	for (let i = 0; i < count; i++) {
		const isAdmin = Math.random() < 0.2 // 20% chance of being admin

		const account = await prisma.account.create({
			data: {
				username: faker.internet.username(),
				password: faker.internet.password(),
				role: isAdmin ? Role.ADMIN : Role.USER,
				...(isAdmin
					? {
							admin_account: {
								create: {},
							},
					  }
					: {
							user_account: {
								create: {
									email: faker.internet.email(),
									phone: faker.phone.number({ style: "international" }),
									gender: randomEnum<Gender>({
										MALE: Gender.MALE,
										FEMALE: Gender.FEMALE,
										OTHER: Gender.OTHER,
									}),
									full_name: faker.person.fullName(),
								},
							},
					  }),
			},
		})
		accounts.push(account)
	}
	return accounts
}

// Generate Brands
async function createBrands(count: number) {
	const brands: Brand[] = []
	for (let i = 0; i < count; i++) {
		const brand = await prisma.brand.create({
			data: {
				name: faker.company.name(),
				description: faker.company.catchPhrase(),
				images: {
					create: [
						{
							url: faker.image.url(),
						},
					],
				},
			},
		})
		brands.push(brand)
	}
	return brands
}

// Generate ProductModels
async function createProductModels(brands: any[], count: number) {
	const productModels: ProductModel[] = []
	for (let i = 0; i < count; i++) {
		const productModel = await prisma.productModel.create({
			data: {
				brand_id: brands[Math.floor(Math.random() * brands.length)].id,
				name: faker.commerce.productName(),
				description: faker.commerce.productDescription(),
				list_price:
					parseInt(faker.commerce.price({ min: 100, max: 1000 })) * 1000,
				date_manufactured: faker.date.past(),
				images: {
					create: [
						{
							url: faker.image.url(),
						},
					],
				},
			},
		})
		productModels.push(productModel)
	}
	return productModels
}

// Generate Products
async function createProducts(productModels: any[], count: number) {
	const products: Product[] = []
	for (let i = 0; i < count; i++) {
		const product = await prisma.product.create({
			data: {
				serial_id: faker.string.alphanumeric(10).toUpperCase(),
				product_model_id:
					productModels[Math.floor(Math.random() * productModels.length)].id,
			},
		})
		products.push(product)
	}
	return products
}

// Main seeding function
async function main() {
	try {
		console.log("Starting to seed database...")

		// Create base data
		const accounts = await createAccounts(10)
		const brands = await createBrands(5)
		const productModels = await createProductModels(brands, 20)
		const products = await createProducts(productModels, 50)

		console.log("Seeding completed successfully!")
	} catch (error) {
		console.error("Error seeding database:", error)
	} finally {
		await prisma.$disconnect()
	}
}

main()
