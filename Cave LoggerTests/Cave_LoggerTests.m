//
//  Cave_LoggerTests.m
//  Cave LoggerTests
//
//  Created by Eoghan Conlon O'Neill on 29/11/2017.
//  Copyright (c) 2017 Eoghan Conlon O'Neill. All rights reserved.
//

#import <Cocoa/Cocoa.h>
#import <XCTest/XCTest.h>

@interface Cave_LoggerTests : XCTestCase

@end

@implementation Cave_LoggerTests

- (void)setUp {
    [super setUp];
    // Put setup code here. This method is called before the invocation of each test method in the class.
}

- (void)tearDown {
    // Put teardown code here. This method is called after the invocation of each test method in the class.
    [super tearDown];
}

- (void)testExample {
    // This is an example of a functional test case.
    XCTAssert(YES, @"Pass");
}

- (void)testPerformanceExample {
    // This is an example of a performance test case.
    [self measureBlock:^{
        // Put the code you want to measure the time of here.
    }];
}

@end
