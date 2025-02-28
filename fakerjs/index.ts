import { faker } from "@faker-js/faker"
import {
	PrismaClient,
	Role,
	Gender,
	Brand,
	Account,
	ProductModel,
	Product,
	PaymentMethod,
	Status,
	RefundMethod,
	Prisma,
} from "@prisma/client"
const prisma = new PrismaClient()

// Helper function to generate random enum values
const randomEnum = <T>(enumObject: { [key: string]: T }): T => {
	const values = Object.values(enumObject)
	return values[Math.floor(Math.random() * values.length)]
}

// Generate Account with UserAccount or AdminAccount
async function createAccounts(count: number) {
	const accountsData: any[] = []

	for (let i = 0; i < count; i++) {
		try {
			const isAdmin = Math.random() < 0.2 // 20% chance of being admin
			const username = faker.internet.username()

			accountsData.push({
				username,
				password: faker.internet.password(),
				role: isAdmin ? Role.ADMIN : Role.USER,
			})
		} catch (error) {
			if (
				!(
					error instanceof Prisma.PrismaClientKnownRequestError &&
					error.code === "P2002"
				)
			) {
				throw error
			}
			// Skip if unique constraint error (P2002)
			continue
		}
	}

	await prisma.account.createMany({
		data: accountsData,
		skipDuplicates: true,
	})

	// Fetch created accounts to get their IDs
	const createdAccounts = await prisma.account.findMany({
		where: {
			username: {
				in: accountsData.map((a) => a.username),
			},
		},
	})

	// Create related admin and user accounts individually since we need the base account IDs
	for (const account of createdAccounts) {
		if (account.role === Role.ADMIN) {
			await prisma.adminAccount.create({
				data: {
					id: account.id,
				},
			})
		} else {
			await prisma.userAccount.create({
				data: {
					id: account.id,
					email: faker.internet.email(),
					phone: faker.phone.number({ style: "international" }),
					gender: randomEnum<Gender>({
						MALE: Gender.MALE,
						FEMALE: Gender.FEMALE,
						OTHER: Gender.OTHER,
					}),
					full_name: faker.person.fullName(),
				},
			})
		}
	}

	return createdAccounts
}

// Generate Brands
async function createBrands(count: number) {
	const brandsData = Array.from({ length: count }, () => ({
		name: faker.company.name(),
		description: faker.company.catchPhrase(),
	}))

	const result = await prisma.brand.createMany({
		data: brandsData,
		skipDuplicates: true,
	})
	return await prisma.brand.findMany({ take: count })
}

// Generate Tags
async function createTags(count: number) {
	const tagsData = Array.from({ length: count }, () => ({
		tag_name: faker.commerce.department(),
		description: faker.commerce.productDescription(),
	}))

	const result = await prisma.tag.createMany({
		data: tagsData,
		skipDuplicates: true,
	})
	return await prisma.tag.findMany({ take: count })
}

// Generate ProductModels
async function createProductModels(
	brands: Brand[],
	tags: any[],
	count: number
) {
	const productModelsData: any[] = []
	const tagOnProductData: any[] = []

	for (let i = 0; i < count; i++) {
		// Randomly assign 1-3 tags to each product model
		const tagCount = Math.floor(Math.random() * 3) + 1
		const selectedTags = faker.helpers.arrayElements(tags, tagCount)
		const modelName = faker.commerce.productName()

		productModelsData.push({
			brand_id: brands[Math.floor(Math.random() * brands.length)].id,
			name: modelName,
			description: faker.commerce.productDescription(),
			list_price: BigInt(
				parseInt(faker.commerce.price({ min: 100, max: 1000 })) * 1000
			),
			date_manufactured: faker.date.past(),
		})

		// We'll need to fetch the created product models to get their IDs for the tag connections
	}

	await prisma.productModel.createMany({
		data: productModelsData,
		skipDuplicates: true,
	})

	// Fetch created product models to get their IDs
	const createdProductModels = await prisma.productModel.findMany({
		take: count,
	})

	// Create tag connections
	for (const productModel of createdProductModels) {
		const tagCount = Math.floor(Math.random() * 3) + 1
		const selectedTags = faker.helpers.arrayElements(tags, tagCount)

		for (const tag of selectedTags) {
			tagOnProductData.push({
				product_model_id: productModel.id,
				tag_name: tag.tag_name,
			})
		}
	}

	if (tagOnProductData.length > 0) {
		await prisma.tagOnProduct.createMany({
			data: tagOnProductData,
			skipDuplicates: true,
		})
	}

	return createdProductModels
}

// Generate Products
async function createProducts(productModels: ProductModel[], count: number) {
	const productsData = Array.from({ length: count }, () => ({
		serial_id: faker.string.alphanumeric(10).toUpperCase(),
		product_model_id:
			productModels[Math.floor(Math.random() * productModels.length)].id,
		sold: Math.random() < 0.3,
	}))

	const result = await prisma.product.createMany({
		data: productsData,
		skipDuplicates: true,
	})
	return await prisma.product.findMany({ take: count })
}

// Generate Addresses
async function createAddresses(userAccounts: any[], count: number) {
	const addressesData = Array.from({ length: count }, () => {
		const userAccount =
			userAccounts[Math.floor(Math.random() * userAccounts.length)]
		return {
			user_id: userAccount.id,
			address: faker.location.streetAddress(),
			city: faker.location.city(),
			province: faker.location.state(),
			country: faker.location.country(),
			postal_code: faker.location.zipCode(),
		}
	})

	await prisma.address.createMany({
		data: addressesData,
		skipDuplicates: true,
	})

	return await prisma.address.findMany({ take: count })
}

// Generate Carts and Items
async function createCarts(userAccounts: any[], productModels: ProductModel[]) {
	const cartsData: any[] = []
	const itemOnCartData: any[] = []

	for (const userAccount of userAccounts) {
		cartsData.push({
			id: userAccount.id,
		})

		// Create cart items for each user
		const itemCount = Math.floor(Math.random() * 5) + 1
		for (let i = 0; i < itemCount; i++) {
			itemOnCartData.push({
				cart_id: userAccount.id,
				product_model_id:
					productModels[Math.floor(Math.random() * productModels.length)].id,
				quantity: BigInt(Math.floor(Math.random() * 5) + 1),
			})
		}
	}

	await prisma.cart.createMany({
		data: cartsData,
		skipDuplicates: true,
	})

	if (itemOnCartData.length > 0) {
		await prisma.itemOnCart.createMany({
			data: itemOnCartData,
			skipDuplicates: true,
		})
	}

	return await prisma.cart.findMany({ take: userAccounts.length })
}

// Generate Sales
async function createSales(
	productModels: ProductModel[],
	tags: any[],
	count: number
) {
	const salesData: any[] = []

	for (let i = 0; i < count; i++) {
		const isByTag = Math.random() < 0.5
		const startDate = faker.date.recent()
		const endDate = faker.date.future({ refDate: startDate })

		salesData.push({
			tag_name: isByTag
				? tags[Math.floor(Math.random() * tags.length)].tag_name
				: null,
			product_model_id: isByTag
				? null
				: productModels[Math.floor(Math.random() * productModels.length)].id,
			date_started: startDate,
			date_ended: endDate,
			quantity: BigInt(Math.floor(Math.random() * 100) + 10),
			used: BigInt(Math.floor(Math.random() * 10)),
			is_active: true,
			discount_percent:
				Math.random() < 0.7 ? BigInt(Math.floor(Math.random() * 50) + 5) : null,
			discount_price:
				Math.random() < 0.3
					? BigInt(Math.floor(Math.random() * 100000) + 10000)
					: null,
		})
	}

	await prisma.sale.createMany({
		data: salesData,
		skipDuplicates: true,
	})

	return await prisma.sale.findMany({ take: count })
}

// Generate Payments and ProductOnPayment
async function createPayments(
	userAccounts: any[],
	products: Product[],
	count: number
) {
	const paymentsData: any[] = []
	const productOnPaymentData: any[] = []
	const productsToUpdate: any[] = []
	const paymentMethods = Object.values(PaymentMethod)
	const statuses = Object.values(Status)

	for (let i = 0; i < count; i++) {
		const userAccount =
			userAccounts[Math.floor(Math.random() * userAccounts.length)]
		const productCount = Math.floor(Math.random() * 3) + 1
		const selectedProducts = faker.helpers.arrayElements(
			products.filter((p) => !p.sold),
			Math.min(productCount, products.filter((p) => !p.sold).length)
		)

		if (selectedProducts.length === 0) continue

		const totalPrice = BigInt(Math.floor(Math.random() * 1000000) + 100000)
		const status = faker.helpers.arrayElement(statuses)

		paymentsData.push({
			user_id: userAccount.id,
			method: faker.helpers.arrayElement(paymentMethods),
			status: status,
			address: faker.location.streetAddress(),
			total: totalPrice,
		})
	}

	// Create payments
	if (paymentsData.length > 0) {
		await prisma.payment.createMany({
			data: paymentsData,
			skipDuplicates: true,
		})
	}

	// Fetch created payments to get their IDs
	const createdPayments = await prisma.payment.findMany({ take: count })

	// Create product on payment records
	for (const payment of createdPayments) {
		const productCount = Math.floor(Math.random() * 3) + 1
		const selectedProducts = faker.helpers.arrayElements(
			products.filter((p) => !p.sold),
			Math.min(productCount, products.filter((p) => !p.sold).length)
		)

		for (const product of selectedProducts) {
			const price = BigInt(Math.floor(Math.random() * 500000) + 50000)
			productOnPaymentData.push({
				payment_id: payment.id,
				product_serial_id: product.serial_id,
				quantity: BigInt(1),
				price: price,
				total_price: price,
			})

			// Track products to update if payment is successful
			if (payment.status === "SUCCESS") {
				productsToUpdate.push(product.id)
			}
		}
	}

	// Create product on payment records
	if (productOnPaymentData.length > 0) {
		await prisma.productOnPayment.createMany({
			data: productOnPaymentData,
			skipDuplicates: true,
		})
	}

	// Update sold status for products with successful payments
	if (productsToUpdate.length > 0) {
		await prisma.product.updateMany({
			where: { id: { in: productsToUpdate } },
			data: { sold: true },
		})
	}

	return createdPayments
}

// Generate Refunds
async function createRefunds(payments: any[], count: number) {
	const refundsData: any[] = []
	const productsToUpdate: any[] = []
	const refundMethods = Object.values(RefundMethod)
	const statuses = Object.values(Status)

	// Only create refunds for successful payments
	const successfulPayments = payments.filter((p) => p.status === "SUCCESS")
	const refundCount = Math.min(count, successfulPayments.length)

	for (let i = 0; i < refundCount; i++) {
		const payment = successfulPayments[i]
		const refundMethod = faker.helpers.arrayElement(refundMethods)
		const status = faker.helpers.arrayElement(statuses)

		refundsData.push({
			payment_id: payment.id,
			method: refundMethod,
			status: status,
			reason: faker.lorem.sentence(),
			address:
				refundMethod === "PICK_UP" ? faker.location.streetAddress() : null,
		})

		// Track products to update if refund is successful
		if (status === "SUCCESS") {
			const productsOnPayment = await prisma.productOnPayment.findMany({
				where: { payment_id: payment.id },
				include: { product: true },
			})

			productsOnPayment.forEach((pop) => {
				productsToUpdate.push(pop.product.id)
			})
		}
	}

	// Create refunds
	if (refundsData.length > 0) {
		await prisma.refund.createMany({
			data: refundsData,
			skipDuplicates: true,
		})
	}

	// Update product sold status for successful refunds
	if (productsToUpdate.length > 0) {
		await prisma.product.updateMany({
			where: { id: { in: productsToUpdate } },
			data: { sold: false },
		})
	}

	return await prisma.refund.findMany({ take: count })
}

// Main seeding function
async function main() {
	try {
		console.log("Starting to seed database...")

		// Create base data
		const accounts = await createAccounts(10)
		const userAccounts = accounts.filter((a) => a.role === "USER")
		const brands = await createBrands(5)
		const tags = await createTags(8)
		const productModels = await createProductModels(brands, tags, 20)
		const products = await createProducts(productModels, 50)

		// Create related data
		const addresses = await createAddresses(userAccounts, 15)
		const carts = await createCarts(userAccounts, productModels)
		const sales = await createSales(productModels, tags, 10)
		const payments = await createPayments(userAccounts, products, 15)
		const refunds = await createRefunds(payments, 5)

		console.log("Seeding completed successfully!")
	} catch (error) {
		console.error("Error seeding database:", error)
	} finally {
		await prisma.$disconnect()
	}
}

main()
